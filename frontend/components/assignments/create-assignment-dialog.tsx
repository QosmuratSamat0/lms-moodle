"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { CalendarIcon, Plus, Loader2 } from "lucide-react";
import { format } from "date-fns";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
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

interface CreateAssignmentDialogProps {
  courseId: string;
  trigger?: React.ReactNode;
}

export function CreateAssignmentDialog({
  courseId,
  trigger,
}: CreateAssignmentDialogProps) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [dueAt, setDueAt] = useState<Date | undefined>();
  const [maxPoints, setMaxPoints] = useState("100");
  const [errors, setErrors] = useState<{ title?: string }>({});
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: () =>
      assignmentService.create({
        course_id: courseId,
        title,
        description: description || undefined,
        due_at: dueAt?.toISOString(),
        max_points: parseInt(maxPoints) || 100,
      }),
    onSuccess: () => {
      toast.success("Assignment created successfully!");
      queryClient.invalidateQueries({
        queryKey: ["course-assignments", courseId],
      });
      queryClient.invalidateQueries({ queryKey: ["assignments"] });
      setOpen(false);
      resetForm();
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to create assignment");
    },
  });

  const resetForm = () => {
    setTitle("");
    setDescription("");
    setDueAt(undefined);
    setMaxPoints("100");
    setErrors({});
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Validate
    const newErrors: { title?: string } = {};
    if (!title || title.length < 3) {
      newErrors.title = "Title must be at least 3 characters";
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    mutation.mutate();
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(v) => {
        setOpen(v);
        if (!v) resetForm();
      }}
    >
      <DialogTrigger asChild>
        {trigger || (
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            New Assignment
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-125">
        <DialogHeader>
          <DialogTitle>Create New Assignment</DialogTitle>
          <DialogDescription>
            Create a new assignment for your students. Set a deadline and
            maximum points.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title *</Label>
            <Input
              id="title"
              placeholder="Assignment title..."
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                if (errors.title) setErrors({});
              }}
            />
            {errors.title && (
              <p className="text-sm text-destructive">{errors.title}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              placeholder="Describe the assignment requirements..."
              className="resize-none min-h-25"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
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
                    disabled={(date) => date < new Date()}
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
              <p className="text-xs text-muted-foreground">Optional deadline</p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="max_points">Max Points</Label>
              <Input
                id="max_points"
                type="number"
                min={0}
                max={1000}
                value={maxPoints}
                onChange={(e) => setMaxPoints(e.target.value)}
              />
              <p className="text-xs text-muted-foreground">0-1000 points</p>
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Create Assignment
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
