"use client";

import { useState, useRef } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  FileText,
  GraduationCap,
  BookOpen,
  Link as LinkIcon,
  File,
  Loader2,
  Plus,
  Upload,
  X,
} from "lucide-react";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import assignmentService from "@/services/assignments";
import uploadService from "@/services/uploads";

export type ContentType =
  | "assignment"
  | "quiz"
  | "lecture"
  | "resource"
  | "link";

interface ContentTypeOption {
  value: ContentType;
  label: string;
  icon: React.ReactNode;
  description: string;
}

const contentTypes: ContentTypeOption[] = [
  {
    value: "assignment",
    label: "Assignment",
    icon: <FileText className="h-5 w-5" />,
    description: "Create a new assignment for students to submit",
  },
  {
    value: "quiz",
    label: "Quiz/Exam",
    icon: <GraduationCap className="h-5 w-5" />,
    description: "Create a quiz or exam (midterm 15%, endterm 15%, etc.)",
  },
  {
    value: "lecture",
    label: "Lecture",
    icon: <BookOpen className="h-5 w-5" />,
    description: "Add lecture materials, slides, or recordings",
  },
  {
    value: "resource",
    label: "Resource/File",
    icon: <File className="h-5 w-5" />,
    description: "Upload syllabus, PDF, or other files",
  },
  {
    value: "link",
    label: "External Link",
    icon: <LinkIcon className="h-5 w-5" />,
    description: "Add a link to YouTube, book, or external resource",
  },
];

interface CreateContentDialogProps {
  courseId: string;
  weekNumber: number | null;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  onSuccess?: () => void;
}

export function CreateContentDialog({
  courseId,
  weekNumber,
  open: controlledOpen,
  onOpenChange: controlledOnOpenChange,
  onSuccess,
}: CreateContentDialogProps) {
  const [internalOpen, setInternalOpen] = useState(false);
  const isControlled = controlledOpen !== undefined;
  const open = isControlled ? controlledOpen : internalOpen;
  const setOpen = isControlled
    ? (v: boolean) => controlledOnOpenChange?.(v)
    : setInternalOpen;

  const [step, setStep] = useState<"select" | "form">("select");
  const [selectedType, setSelectedType] = useState<ContentType | null>(null);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [dueDate, setDueDate] = useState("");
  const [weightPercentage, setWeightPercentage] = useState("");
  const [url, setUrl] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const queryClient = useQueryClient();

  const actualWeek = weekNumber ?? 0;

  const resetForm = () => {
    setStep("select");
    setSelectedType(null);
    setTitle("");
    setDescription("");
    setDueDate("");
    setWeightPercentage("");
    setUrl("");
    setSelectedFile(null);
  };

  const createAssignmentMutation = useMutation({
    mutationFn: (data: {
      course_id: string;
      title: string;
      description?: string;
      due_at?: string;
      max_points: number;
      weight_percentage?: number;
      file_url?: string;
    }) => assignmentService.create(data),
    onSuccess: async () => {
      toast.success("Content created successfully!");
      queryClient.invalidateQueries({
        queryKey: ["course-assignments", courseId],
      });
      setOpen(false);
      resetForm();
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to create content");
    },
  });

  const handleSelectType = (type: ContentType) => {
    setSelectedType(type);
    setStep("form");
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      // Check file size (max 50MB)
      if (file.size > 50 * 1024 * 1024) {
        toast.error("File size must be less than 50MB");
        return;
      }
      setSelectedFile(file);
    }
  };

  const handleSubmit = async () => {
    if (!title.trim()) {
      toast.error("Please enter a title");
      return;
    }

    const titleWithWeek =
      actualWeek > 0 ? `Week ${actualWeek}: ${title}` : title;
    const weight = parseFloat(weightPercentage) || 0;

    // Build clean description — no category prefixes, no [FILE:] tags
    let fullDescription = description || "";

    if (selectedType === "quiz") {
      fullDescription = `[QUIZ] ${description || ""}`;
    } else if (selectedType === "lecture") {
      fullDescription = `[LECTURE] ${description || ""}`;
    } else if (selectedType === "resource") {
      fullDescription = `[RESOURCE] ${description || ""}`;
    } else if (selectedType === "link") {
      fullDescription = `[LINK] ${url}\n${description || ""}`;
    }

    // Upload file first if present, get URL
    let fileUrl: string | undefined;
    if (selectedFile && selectedType !== "link") {
      try {
        setIsUploading(true);
        const uploaded = await uploadService.upload(selectedFile, "assignment");
        fileUrl = uploaded.secure_url || uploaded.url;
      } catch {
        toast.warning("File upload failed, creating content without file");
      } finally {
        setIsUploading(false);
      }
    }

    if (selectedType === "assignment" || selectedType === "quiz") {
      createAssignmentMutation.mutate({
        course_id: courseId,
        title: titleWithWeek,
        description: fullDescription,
        due_at: dueDate || undefined,
        max_points: 100,
        weight_percentage: weight,
        file_url: fileUrl,
      });
    } else {
      // For lectures, resources, links - create as assignment with 0 weight
      createAssignmentMutation.mutate({
        course_id: courseId,
        title: titleWithWeek,
        description: fullDescription,
        max_points: 0,
        file_url: fileUrl,
      });
    }
  };

  const isPending = createAssignmentMutation.isPending || isUploading;

  const getPlaceholder = () => {
    switch (selectedType) {
      case "assignment":
        return "e.g., Assignment 1, Homework 3, Project 2";
      case "quiz":
        return "e.g., Quiz 1, Midterm Exam, Final Quiz";
      case "lecture":
        return "e.g., Lecture 1, Introduction to Programming";
      case "resource":
        return "e.g., Syllabus, Course Notes, Reference PDF";
      case "link":
        return "e.g., YouTube Tutorial, Textbook Chapter 1";
      default:
        return "Enter title";
    }
  };

  const weekLabel = actualWeek > 0 ? `Week ${actualWeek}` : "Course Materials";

  return (
    <Dialog
      open={open}
      onOpenChange={(value) => {
        setOpen(value);
        if (!value) resetForm();
      }}
    >
      {!isControlled && (
        <DialogTrigger asChild>
          <Button
            variant="ghost"
            size="sm"
            className="gap-1 text-primary hover:text-primary"
          >
            <Plus className="h-4 w-4" />
            Add content
          </Button>
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>
            {step === "select"
              ? "Add Content"
              : `Add ${contentTypes.find((t) => t.value === selectedType)?.label}`}
          </DialogTitle>
          <DialogDescription>
            {step === "select"
              ? `Add content to ${weekLabel}`
              : `Create a new ${selectedType} for ${weekLabel}`}
          </DialogDescription>
        </DialogHeader>

        {step === "select" ? (
          <div className="grid gap-3 py-4">
            {contentTypes.map((type) => (
              <button
                key={type.value}
                onClick={() => handleSelectType(type.value)}
                className="flex items-center gap-4 p-4 rounded-lg border hover:bg-muted/50 hover:border-primary/50 transition-colors text-left"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                  {type.icon}
                </div>
                <div className="flex-1">
                  <p className="font-medium">{type.label}</p>
                  <p className="text-sm text-muted-foreground">
                    {type.description}
                  </p>
                </div>
              </button>
            ))}
          </div>
        ) : (
          <div className="grid gap-4 py-4 max-h-[60vh] overflow-y-auto">
            <div className="grid gap-2">
              <Label htmlFor="title">Title *</Label>
              <Input
                id="title"
                placeholder={getPlaceholder()}
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                disabled={isPending}
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                placeholder="Enter description (optional)"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={3}
                disabled={isPending}
              />
            </div>

            {selectedType === "link" && (
              <div className="grid gap-2">
                <Label htmlFor="url">URL *</Label>
                <Input
                  id="url"
                  type="url"
                  placeholder="https://youtube.com/watch?v=..."
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  disabled={isPending}
                />
              </div>
            )}

            {/* File Upload for assignments, lectures, resources */}
            {(selectedType === "assignment" ||
              selectedType === "lecture" ||
              selectedType === "resource" ||
              selectedType === "quiz") && (
              <div className="grid gap-2">
                <Label>Attachment (optional)</Label>
                <input
                  ref={fileInputRef}
                  type="file"
                  className="hidden"
                  onChange={handleFileChange}
                  disabled={isPending}
                />
                {selectedFile ? (
                  <div className="flex items-center gap-2 p-3 border rounded-lg bg-muted/50">
                    <File className="h-5 w-5 text-primary" />
                    <span className="flex-1 truncate text-sm">
                      {selectedFile.name}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {(selectedFile.size / 1024).toFixed(1)} KB
                    </span>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => setSelectedFile(null)}
                      disabled={isPending}
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
                    disabled={isPending}
                  >
                    <Upload className="h-4 w-4 mr-2" />
                    Upload File (PDF, DOC, ZIP, etc.)
                  </Button>
                )}
                <p className="text-xs text-muted-foreground">
                  Max file size: 50MB. Students can download this file.
                </p>
              </div>
            )}

            {(selectedType === "assignment" || selectedType === "quiz") && (
              <>
                <div className="grid grid-cols-2 gap-4">
                  <div className="grid gap-2">
                    <Label htmlFor="due-date">Due Date</Label>
                    <Input
                      id="due-date"
                      type="datetime-local"
                      value={dueDate}
                      onChange={(e) => setDueDate(e.target.value)}
                      disabled={isPending}
                    />
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="weight">Weight (%)</Label>
                    <div className="relative">
                      <Input
                        id="weight"
                        type="number"
                        min="0"
                        max="100"
                        step="1"
                        placeholder="100"
                        value={weightPercentage}
                        onChange={(e) => setWeightPercentage(e.target.value)}
                        disabled={isPending}
                      />
                    </div>
                  </div>
                </div>
              </>
            )}
          </div>
        )}

        <DialogFooter>
          {step === "form" && (
            <>
              <Button
                variant="outline"
                onClick={() => setStep("select")}
                disabled={isPending}
              >
                Back
              </Button>
              <Button onClick={handleSubmit} disabled={isPending}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Create
              </Button>
            </>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
