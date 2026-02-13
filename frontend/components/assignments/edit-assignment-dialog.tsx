"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { CalendarIcon, Loader2, Pencil, Trash2 } from "lucide-react";
import { format } from "date-fns";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { cn } from "@/lib/utils";
import assignmentService from "@/services/assignments";
import type { Assignment } from "@/types/assignment";

// Clean description (remove legacy category tags)
const cleanDescription = (description?: string): string => {
  if (!description) return "";
  return description
    .replace(/\[REGISTER_MIDTERM\]\s*/i, "")
    .replace(/\[REGISTER_ENDTERM\]\s*/i, "")
    .replace(/\[FINAL\]\s*/i, "")
    .replace(/\[BONUS\]\s*/i, "")
    .replace(/\[QUIZ\]\s*/i, "")
    .replace(/\n\[FILE:[^\]]*\]/g, "")
    .trim();
};

interface EditAssignmentDialogProps {
  assignment: Assignment;
  courseId: string;
  trigger?: React.ReactNode;
}

export function EditAssignmentDialog({
  assignment,
  courseId,
  trigger,
}: EditAssignmentDialogProps) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState(assignment.title);
  const [description, setDescription] = useState(
    cleanDescription(assignment.description),
  );
  const [dueAt, setDueAt] = useState<Date | undefined>(
    assignment.due_at ? new Date(assignment.due_at) : undefined,
  );
  const [weightPercentage, setWeightPercentage] = useState(
    String(assignment.max_points || ""),
  );
  const [errors, setErrors] = useState<{ title?: string; weight?: string }>({});
  const queryClient = useQueryClient();

  const weight = parseFloat(weightPercentage) || 0;

  const handleOpenChange = (isOpen: boolean) => {
    setOpen(isOpen);
    if (isOpen) {
      // Reset form with current assignment data when opening
      setTitle(assignment.title);
      setDescription(cleanDescription(assignment.description));
      setDueAt(assignment.due_at ? new Date(assignment.due_at) : undefined);
      setWeightPercentage(String(assignment.max_points || ""));
      setErrors({});
    }
  };

  const updateMutation = useMutation({
    mutationFn: () =>
      assignmentService.update(assignment.id, {
        title,
        description: description || undefined,
        due_at: dueAt?.toISOString(),
        max_points: weight,
        weight_percentage: weight,
      }),
    onSuccess: () => {
      toast.success("Assignment updated successfully!");
      queryClient.invalidateQueries({
        queryKey: ["course-assignments", courseId],
      });
      queryClient.invalidateQueries({ queryKey: ["my-course-assignments"] });
      queryClient.invalidateQueries({ queryKey: ["my-deadline-assignments"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard-assignments"] });
      queryClient.invalidateQueries({ queryKey: ["assignments"] });
      queryClient.invalidateQueries({
        queryKey: ["assignment", assignment.id],
      });
      setOpen(false);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to update assignment");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: () => assignmentService.delete(assignment.id),
    onSuccess: () => {
      toast.success("Assignment deleted successfully!");
      queryClient.invalidateQueries({
        queryKey: ["course-assignments", courseId],
      });
      queryClient.invalidateQueries({ queryKey: ["my-course-assignments"] });
      queryClient.invalidateQueries({ queryKey: ["my-deadline-assignments"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard-assignments"] });
      queryClient.invalidateQueries({ queryKey: ["assignments"] });
      setOpen(false);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to delete assignment");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const newErrors: { title?: string; weight?: string } = {};
    if (!title || title.length < 3) {
      newErrors.title = "Title must be at least 3 characters";
    }
    if (weight <= 0 || weight > 100) {
      newErrors.weight = "Weight must be between 0.1% and 100%";
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    updateMutation.mutate();
  };

  const isPending = updateMutation.isPending || deleteMutation.isPending;

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        {trigger || (
          <Button variant="outline" size="sm">
            <Pencil className="mr-2 h-4 w-4" />
            Edit
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Edit Assignment</DialogTitle>
          <DialogDescription>
            Update the assignment details or delete this assignment.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="edit-title">Title *</Label>
            <Input
              id="edit-title"
              placeholder="Assignment title..."
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                if (errors.title) setErrors({});
              }}
              disabled={isPending}
            />
            {errors.title && (
              <p className="text-sm text-destructive">{errors.title}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="edit-description">Description</Label>
            <Textarea
              id="edit-description"
              placeholder="Describe the assignment requirements..."
              className="resize-none min-h-20"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={isPending}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Due Date</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    className={cn(
                      "w-full pl-3 text-left font-normal",
                      !dueAt && "text-muted-foreground",
                    )}
                    disabled={isPending}
                  >
                    {dueAt ? format(dueAt, "PPP") : <span>Pick a date</span>}
                    <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={dueAt}
                    onSelect={setDueAt}
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
              <p className="text-xs text-muted-foreground">Optional deadline</p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="edit-weight">Max Points</Label>
              <Input
                id="edit-weight"
                type="number"
                min={0}
                max={1000}
                step={1}
                placeholder="100"
                value={weightPercentage}
                onChange={(e) => {
                  setWeightPercentage(e.target.value);
                  if (errors.weight)
                    setErrors((prev) => ({ ...prev, weight: undefined }));
                }}
                disabled={isPending}
              />
              {errors.weight && (
                <p className="text-sm text-destructive">{errors.weight}</p>
              )}
            </div>
          </div>

          <div className="flex justify-between pt-4">
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button
                  type="button"
                  variant="destructive"
                  disabled={isPending}
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Delete Assignment</AlertDialogTitle>
                  <AlertDialogDescription>
                    Are you sure you want to delete this assignment? This action
                    cannot be undone. All submissions will also be deleted.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction
                    onClick={() => deleteMutation.mutate()}
                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                  >
                    {deleteMutation.isPending && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Delete
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>

            <div className="flex gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
                disabled={isPending}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isPending}>
                {updateMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Save Changes
              </Button>
            </div>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
