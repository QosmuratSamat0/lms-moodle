"use client";

import { useState, useRef } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  CalendarIcon,
  Plus,
  Loader2,
  Percent,
  Upload,
  X,
  FileText,
} from "lucide-react";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";
import assignmentService from "@/services/assignments";
import uploadService from "@/services/uploads";
import { GradingCategory } from "@/types/assignment";

const CATEGORY_LABELS: Record<GradingCategory, string> = {
  register_midterm: "Register Midterm (30%)",
  register_endterm: "Register Endterm (30%)",
  final: "Final Exam (40%)",
  bonus: "Bonus (Extra Credit)",
};

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
  const [weightPercentage, setWeightPercentage] = useState("");
  const [gradingCategory, setGradingCategory] =
    useState<GradingCategory>("register_midterm");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadedFileUrl, setUploadedFileUrl] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [errors, setErrors] = useState<{ title?: string; weight?: string }>({});
  const fileInputRef = useRef<HTMLInputElement>(null);
  const queryClient = useQueryClient();

  const weight = parseFloat(weightPercentage) || 0;

  const mutation = useMutation({
    mutationFn: async () => {
      // Upload file first if selected and not yet uploaded
      let fileUrl = uploadedFileUrl;
      if (selectedFile && !fileUrl) {
        setIsUploading(true);
        try {
          const uploaded = await uploadService.upload(
            selectedFile,
            "assignment",
          );
          fileUrl = uploaded.secure_url || uploaded.url;
          setUploadedFileUrl(fileUrl);
        } catch {
          throw new Error("Failed to upload file");
        } finally {
          setIsUploading(false);
        }
      }

      return assignmentService.create({
        course_id: courseId,
        title,
        description: description || undefined,
        due_at: dueAt?.toISOString(),
        max_points: weight,
        weight_percentage: weight,
        grading_category: gradingCategory,
        file_url: fileUrl || undefined,
      });
    },
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
    setWeightPercentage("");
    setGradingCategory("register_midterm");
    setSelectedFile(null);
    setUploadedFileUrl(null);
    setErrors({});
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 50 * 1024 * 1024) {
        toast.error("File size must be less than 50MB");
        return;
      }
      setSelectedFile(file);
      setUploadedFileUrl(null);
    }
  };

  const removeFile = () => {
    setSelectedFile(null);
    setUploadedFileUrl(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

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
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Create New Assignment</DialogTitle>
          <DialogDescription>
            Set the grading category and weight for this assignment.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title *</Label>
            <Input
              id="title"
              placeholder="e.g., Assignment 1, Homework 3, Project 2"
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
              placeholder="Enter description (optional)"
              className="resize-none min-h-20"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>

          {/* File Attachment */}
          <div className="space-y-2">
            <Label>Attachment (optional)</Label>
            <input
              ref={fileInputRef}
              type="file"
              className="hidden"
              onChange={handleFileSelect}
              accept=".pdf,.doc,.docx,.zip,.rar,.ppt,.pptx,.xls,.xlsx,.txt,.png,.jpg,.jpeg"
            />
            {selectedFile ? (
              <div className="flex items-center gap-2 p-3 rounded-lg border bg-muted/50">
                <FileText className="h-5 w-5 text-primary shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">
                    {selectedFile.name}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                  </p>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={removeFile}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            ) : (
              <Button
                type="button"
                variant="outline"
                className="w-full"
                onClick={() => fileInputRef.current?.click()}
              >
                <Upload className="mr-2 h-4 w-4" />
                Upload File (PDF, DOC, ZIP, etc.)
              </Button>
            )}
            <p className="text-xs text-muted-foreground">
              Max file size: 50MB. Students can download this file.
            </p>
          </div>

          {/* Grading Category */}
          <div className="space-y-2">
            <Label>Grading Category *</Label>
            <Select
              value={gradingCategory}
              onValueChange={(v) => setGradingCategory(v as GradingCategory)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="register_midterm">
                  {CATEGORY_LABELS.register_midterm}
                </SelectItem>
                <SelectItem value="register_endterm">
                  {CATEGORY_LABELS.register_endterm}
                </SelectItem>
                <SelectItem value="final">{CATEGORY_LABELS.final}</SelectItem>
                <SelectItem value="bonus">{CATEGORY_LABELS.bonus}</SelectItem>
              </SelectContent>
            </Select>
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
            </div>

            <div className="space-y-2">
              <Label htmlFor="weight">Weight *</Label>
              <div className="relative">
                <Input
                  id="weight"
                  type="number"
                  min={0.1}
                  max={100}
                  step={0.1}
                  placeholder="e.g., 5"
                  value={weightPercentage}
                  onChange={(e) => {
                    setWeightPercentage(e.target.value);
                    if (errors.weight)
                      setErrors((prev) => ({ ...prev, weight: undefined }));
                  }}
                  className="pr-8"
                />
                <Percent className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              </div>
              {errors.weight && (
                <p className="text-sm text-destructive">{errors.weight}</p>
              )}
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
            <Button type="submit" disabled={mutation.isPending || isUploading}>
              {(mutation.isPending || isUploading) && (
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
