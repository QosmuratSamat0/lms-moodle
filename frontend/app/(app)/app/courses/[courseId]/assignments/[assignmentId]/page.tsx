"use client";

import { useState, useRef, useMemo, useEffect } from "react";
import { motion } from "framer-motion";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import Link from "next/link";
import { toast } from "sonner";
import {
  ArrowLeft,
  Upload,
  FileText,
  Clock,
  CheckCircle2,
  XCircle,
  Loader2,
  AlertTriangle,
  Award,
  Calendar,
  User,
  Download,
  MessageSquare,
  Pencil,
  Trash2,
  Star,
  UsersRound,
} from "lucide-react";

import { PageHeader } from "@/components/common/page-header";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Skeleton } from "@/components/ui/skeleton";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAuthStore } from "@/store/auth-store";
import {
  formatDate,
  formatDateTime,
  isOverdue,
  getSubmissionStatusColor,
  getFullName,
  getTimeRemaining,
  getSubmissionTimeInfo,
  cn,
} from "@/lib/helpers";
import assignmentService from "@/services/assignments";
import submissionService from "@/services/submissions";
import uploadService from "@/services/uploads";
import groupService from "@/services/groups";
import type { Submission } from "@/types/submission";

// Clean description by removing legacy tags
const cleanAssignmentDescription = (desc?: string): string => {
  if (!desc) return "";
  return desc
    .replace(/\[REGISTER_MIDTERM\]\s*/gi, "")
    .replace(/\[REGISTER_ENDTERM\]\s*/gi, "")
    .replace(/\[FINAL\]\s*/gi, "")
    .replace(/\[BONUS\]\s*/gi, "")
    .replace(/\[QUIZ\]\s*/gi, "")
    .replace(/\[LECTURE\]\s*/gi, "")
    .replace(/\[RESOURCE\]\s*/gi, "")
    .replace(/\[LINK\]\s*/gi, "")
    .replace(/\n?\[FILE:[^\]]*\]/g, "")
    .trim();
};

export default function AssignmentDetailPage() {
  const params = useParams();
  const courseId = params.courseId as string;
  const assignmentId = params.assignmentId as string;
  const queryClient = useQueryClient();
  const { isTeacher, isStudent, user } = useAuthStore();

  const [submissionText, setSubmissionText] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isEditMode, setIsEditMode] = useState(false);

  // Selected group for teachers (to view specific group's submissions)
  const [selectedGroupId, setSelectedGroupId] = useState<string | null>(null);

  // Ref for the edit file input to trigger click programmatically
  const editFileInputRef = useRef<HTMLInputElement>(null);

  // Teacher grading state
  const [gradingSubmission, setGradingSubmission] = useState<Submission | null>(
    null,
  );
  const [gradeScore, setGradeScore] = useState<string>("");
  const [gradeFeedback, setGradeFeedback] = useState<string>("");

  // Fetch teacher's group assignments for this course
  const { data: teacherAssignments } = useQuery({
    queryKey: ["teacher-group-assignments"],
    queryFn: () => groupService.getMyAssignments({ limit: 100 }),
    enabled: isTeacher(),
  });

  // Filter assignments for this course only and get unique groups
  const courseGroups = useMemo(() => {
    if (!teacherAssignments?.assignments) return [];
    return teacherAssignments.assignments
      .filter((a) => a.course_id === courseId)
      .sort((a, b) => a.group_code.localeCompare(b.group_code));
  }, [teacherAssignments, courseId]);

  // Set default selected group when data loads
  useEffect(() => {
    if (courseGroups.length > 0 && !selectedGroupId) {
      setSelectedGroupId(courseGroups[0].group_id);
    }
  }, [courseGroups, selectedGroupId]);

  const { data: assignment, isLoading: assignmentLoading } = useQuery({
    queryKey: ["assignment", assignmentId],
    queryFn: () => assignmentService.get(assignmentId),
  });

  const { data: submissionsData, isLoading: submissionsLoading } = useQuery({
    queryKey: ["assignment-submissions", assignmentId, selectedGroupId],
    queryFn: () =>
      submissionService.getByAssignment(assignmentId, {
        group_id: selectedGroupId || undefined,
      }),
    enabled: isTeacher(),
  });

  const { data: mySubmissions } = useQuery({
    queryKey: ["my-submissions", user?.id],
    queryFn: () => submissionService.getMySubmissions(user!.id),
    enabled: isStudent() && !!user?.id,
  });

  const mySubmission = mySubmissions?.submissions?.find(
    (s) => s.assignment_id === assignmentId,
  );

  const [isUploading, setIsUploading] = useState(false);

  const submitMutation = useMutation({
    mutationFn: async () => {
      let fileUrl: string | undefined;

      // Upload file to Cloudinary first if there's a file
      if (selectedFile) {
        setIsUploading(true);
        try {
          const uploadResult = await uploadService.upload(
            selectedFile,
            "submission",
          );
          // Use secure_url from Cloudinary (or url as fallback)
          fileUrl = uploadResult.secure_url || uploadResult.url;
          console.log("File uploaded to Cloudinary:", fileUrl);
        } catch (error) {
          setIsUploading(false);
          console.error("Upload error:", error);
          throw new Error("Failed to upload file. Please try again.");
        }
        setIsUploading(false);
      }

      // Now create the submission with the Cloudinary URL
      return submissionService.create({
        assignment_id: assignmentId,
        content_text: submissionText || undefined,
        file_url: fileUrl,
      });
    },
    onSuccess: () => {
      toast.success("Submission uploaded successfully!");
      setSubmissionText("");
      setSelectedFile(null);
      queryClient.invalidateQueries({ queryKey: ["my-submissions"] });
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to submit");
    },
  });

  const updateMutation = useMutation({
    mutationFn: async (data: { content_text?: string }) => {
      let fileUrl: string | undefined;

      // Upload new file if selected
      if (selectedFile) {
        setIsUploading(true);
        try {
          const uploadResult = await uploadService.upload(
            selectedFile,
            "submission",
          );
          fileUrl = uploadResult.secure_url || uploadResult.url;
          console.log("File uploaded to Cloudinary:", fileUrl);
        } catch (error) {
          setIsUploading(false);
          console.error("Upload error:", error);
          throw new Error("Failed to upload file. Please try again.");
        }
        setIsUploading(false);
      }

      return submissionService.update(mySubmission?.id || "", {
        content_text: data.content_text,
        file_url: fileUrl,
      });
    },
    onSuccess: () => {
      toast.success("Submission updated successfully!");
      setSubmissionText("");
      setSelectedFile(null);
      setIsEditMode(false);
      queryClient.invalidateQueries({ queryKey: ["my-submissions"] });
      queryClient.invalidateQueries({ queryKey: ["assignment", assignmentId] });
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to update submission");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: () => submissionService.delete(mySubmission?.id || ""),
    onSuccess: () => {
      toast.success("Submission deleted successfully!");
      queryClient.invalidateQueries({ queryKey: ["my-submissions"] });
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to delete submission");
    },
  });

  // Teacher grade mutation
  const gradeMutation = useMutation({
    mutationFn: (data: {
      submission_id: string;
      score: number;
      feedback?: string;
    }) => submissionService.gradeSubmission(data),
    onSuccess: () => {
      toast.success("Submission graded successfully!");
      setGradingSubmission(null);
      setGradeScore("");
      setGradeFeedback("");
      queryClient.invalidateQueries({
        queryKey: ["assignment-submissions", assignmentId],
      });
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to grade submission");
    },
  });

  const handleSubmit = () => {
    if (!submissionText && !selectedFile) {
      toast.error("Please provide either text or a file");
      return;
    }
    submitMutation.mutate();
  };

  // Handle edit button - enter edit mode
  const handleEditClick = () => {
    if (mySubmission) {
      setSubmissionText(mySubmission.content_text || "");
    }
    setIsEditMode(true);
  };

  // Handle edit file selection
  const handleEditFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
    }
    // Reset the input so the same file can be selected again
    if (e.target) {
      e.target.value = "";
    }
  };

  // Handle update submission
  const handleUpdateSubmit = () => {
    if (!submissionText && !selectedFile) {
      toast.error("Please provide either text or a file");
      return;
    }
    updateMutation.mutate({ content_text: submissionText });
  };

  // Cancel edit mode
  const handleCancelEdit = () => {
    setIsEditMode(false);
    setSelectedFile(null);
    setSubmissionText("");
  };

  // Handle grade submission
  const handleGradeSubmit = () => {
    if (!gradingSubmission) return;

    const score = Math.round(parseFloat(gradeScore));
    if (isNaN(score) || score < 0 || score > (assignment?.max_points || 100)) {
      toast.error(
        `Score must be between 0 and ${assignment?.max_points || 100}`,
      );
      return;
    }

    gradeMutation.mutate({
      submission_id: gradingSubmission.id,
      score,
      feedback: gradeFeedback || undefined,
    });
  };

  // Open grading dialog
  const openGradingDialog = (submission: Submission) => {
    setGradingSubmission(submission);
    setGradeScore(submission.grade_score?.toString() || "");
    setGradeFeedback(submission.grade_feedback || "");
  };

  if (assignmentLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-10 w-1/2" />
        <Skeleton className="h-[400px] w-full" />
      </div>
    );
  }

  if (!assignment) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <h2 className="text-xl font-semibold">Assignment not found</h2>
        <Button asChild className="mt-4">
          <Link href={`/app/courses/${courseId}`}>Back to course</Link>
        </Button>
      </div>
    );
  }

  const overdue = assignment.due_at && isOverdue(assignment.due_at);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href={`/app/courses/${courseId}`}>
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <PageHeader
          title={assignment.title}
          description={assignment.course_title}
        />
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Assignment Details */}
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Assignment Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h4 className="font-medium mb-2">Description</h4>
                <p className="text-muted-foreground">
                  {cleanAssignmentDescription(assignment.description) ||
                    "No description provided."}
                </p>
              </div>
              {assignment.file_url && (
                <div>
                  <h4 className="font-medium mb-2">Attachment</h4>
                  <Button variant="outline" size="sm" asChild>
                    <a
                      href={assignment.file_url}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <Download className="h-4 w-4 mr-2" />
                      Download File
                    </a>
                  </Button>
                </div>
              )}
              <Separator />
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <h4 className="font-medium mb-1">Due Date</h4>
                  <div className="flex items-center gap-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <span className={overdue ? "text-destructive" : ""}>
                      {assignment.due_at
                        ? formatDateTime(assignment.due_at)
                        : "No deadline"}
                    </span>
                  </div>
                </div>
                <div>
                  <h4 className="font-medium mb-1">Max Points</h4>
                  <span>{assignment.max_points} pts</span>
                </div>
                <div>
                  <h4 className="font-medium mb-1">Instructor</h4>
                  <span>
                    {getFullName(
                      assignment.teacher_first_name,
                      assignment.teacher_last_name,
                    )}
                  </span>
                </div>
                <div>
                  <h4 className="font-medium mb-1">Created</h4>
                  <span>{formatDate(assignment.created_at)}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Student: Submission Form */}
          {isStudent() && !mySubmission && (
            <Card>
              <CardHeader>
                <CardTitle>Submit Assignment</CardTitle>
                <CardDescription>
                  Upload your work for this assignment
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="submission-text">Text Submission</Label>
                  <Textarea
                    id="submission-text"
                    placeholder="Enter your submission text here..."
                    value={submissionText}
                    onChange={(e) => setSubmissionText(e.target.value)}
                    rows={6}
                    disabled={submitMutation.isPending}
                  />
                </div>
                <div className="space-y-2">
                  <Label>File Upload</Label>
                  <div className="border-2 border-dashed rounded-lg p-8 text-center">
                    <input
                      type="file"
                      id="file-upload"
                      className="hidden"
                      onChange={(e) =>
                        setSelectedFile(e.target.files?.[0] || null)
                      }
                      disabled={submitMutation.isPending}
                    />
                    <label
                      htmlFor="file-upload"
                      className="cursor-pointer flex flex-col items-center gap-2"
                    >
                      <Upload className="h-8 w-8 text-muted-foreground" />
                      {selectedFile ? (
                        <span className="text-sm font-medium">
                          {selectedFile.name}
                        </span>
                      ) : (
                        <>
                          <span className="text-sm font-medium">
                            Click to upload
                          </span>
                          <span className="text-xs text-muted-foreground">
                            PDF, DOC, or ZIP (max 10MB)
                          </span>
                        </>
                      )}
                    </label>
                  </div>
                </div>
                <Button
                  onClick={handleSubmit}
                  disabled={
                    submitMutation.isPending ||
                    isUploading ||
                    (!submissionText && !selectedFile)
                  }
                  className="w-full"
                >
                  {(submitMutation.isPending || isUploading) && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  {isUploading
                    ? "Uploading file..."
                    : submitMutation.isPending
                      ? "Submitting..."
                      : "Submit Assignment"}
                </Button>
              </CardContent>
            </Card>
          )}

          {/* Student: View Submission - Moodle Style */}
          {isStudent() && mySubmission && (
            <Card>
              <CardHeader>
                <CardTitle>Submission status</CardTitle>
              </CardHeader>
              <CardContent className="p-0">
                {/* Submission Status Table */}
                <div className="divide-y">
                  {/* Submission status row */}
                  <div
                    className={cn(
                      "grid grid-cols-3 gap-4 p-4",
                      mySubmission.status === "graded"
                        ? "bg-green-50 dark:bg-green-900/20"
                        : "bg-yellow-50 dark:bg-yellow-900/20",
                    )}
                  >
                    <div className="font-medium">Submission status</div>
                    <div className="col-span-2">
                      {mySubmission.status === "graded"
                        ? "Submitted for grading"
                        : "Submitted for grading"}
                    </div>
                  </div>

                  {/* Grading status row */}
                  <div
                    className={cn(
                      "grid grid-cols-3 gap-4 p-4",
                      mySubmission.status === "graded"
                        ? "bg-green-50 dark:bg-green-900/20"
                        : "bg-yellow-50 dark:bg-yellow-900/20",
                    )}
                  >
                    <div className="font-medium">Grading status</div>
                    <div className="col-span-2">
                      {mySubmission.status === "graded" ? (
                        <span className="text-green-600 dark:text-green-400 font-medium">
                          Graded
                        </span>
                      ) : (
                        <span className="text-yellow-600 dark:text-yellow-400">
                          Not graded
                        </span>
                      )}
                    </div>
                  </div>

                  {/* Time remaining row */}
                  <div
                    className={cn(
                      "grid grid-cols-3 gap-4 p-4",
                      assignment.due_at &&
                        new Date(mySubmission.submitted_at) >
                          new Date(assignment.due_at)
                        ? "bg-red-50 dark:bg-red-900/20"
                        : "bg-green-50 dark:bg-green-900/20",
                    )}
                  >
                    <div className="font-medium">Time remaining</div>
                    <div
                      className={cn(
                        "col-span-2",
                        assignment.due_at &&
                          new Date(mySubmission.submitted_at) >
                            new Date(assignment.due_at)
                          ? "text-red-600 dark:text-red-400"
                          : "text-green-600 dark:text-green-400",
                      )}
                    >
                      {getSubmissionTimeInfo(
                        mySubmission.submitted_at,
                        assignment.due_at,
                      )}
                    </div>
                  </div>

                  {/* Last modified row */}
                  <div className="grid grid-cols-3 gap-4 p-4 bg-muted/30">
                    <div className="font-medium">Last modified</div>
                    <div className="col-span-2">
                      {formatDateTime(mySubmission.submitted_at)}
                    </div>
                  </div>

                  {/* File submissions row */}
                  {mySubmission.file_url && (
                    <div className="grid grid-cols-3 gap-4 p-4 bg-muted/30">
                      <div className="font-medium">File submissions</div>
                      <div className="col-span-2">
                        <div className="flex items-center gap-3 p-3 bg-background rounded-lg border">
                          <FileText className="h-8 w-8 text-blue-500" />
                          <div className="flex-1 min-w-0">
                            <p className="font-medium truncate">
                              {mySubmission.file_name || "Submission file"}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              Uploaded{" "}
                              {formatDateTime(mySubmission.submitted_at)}
                            </p>
                          </div>
                          <div className="flex items-center gap-2">
                            <Button variant="outline" size="sm" asChild>
                              <a
                                href={mySubmission.file_url}
                                target="_blank"
                                rel="noopener noreferrer"
                              >
                                <Download className="h-4 w-4 mr-1" />
                                Download
                              </a>
                            </Button>
                          </div>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Text submission row */}
                  {mySubmission.content_text && (
                    <div className="grid grid-cols-3 gap-4 p-4 bg-muted/30">
                      <div className="font-medium">Online text</div>
                      <div className="col-span-2">
                        <p className="text-sm bg-background p-3 rounded border whitespace-pre-wrap">
                          {mySubmission.content_text}
                        </p>
                      </div>
                    </div>
                  )}

                  {/* Actions: Edit and Delete (Only if not graded) */}
                  {mySubmission.status !== "graded" && !isEditMode && (
                    <div className="grid grid-cols-3 gap-4 p-4 bg-muted/30">
                      <div className="font-medium">Actions</div>
                      <div className="col-span-2 flex gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={handleEditClick}
                        >
                          <Pencil className="h-4 w-4 mr-1" />
                          Edit submission
                        </Button>
                        <Button
                          variant="destructive"
                          size="sm"
                          disabled={deleteMutation.isPending}
                          onClick={() => {
                            if (
                              confirm(
                                "Are you sure you want to delete your submission? This action cannot be undone.",
                              )
                            ) {
                              deleteMutation.mutate();
                            }
                          }}
                        >
                          {deleteMutation.isPending ? (
                            <Loader2 className="h-4 w-4 mr-1 animate-spin" />
                          ) : (
                            <Trash2 className="h-4 w-4 mr-1" />
                          )}
                          Delete
                        </Button>
                      </div>
                    </div>
                  )}

                  {/* Edit Mode Form */}
                  {mySubmission.status !== "graded" && isEditMode && (
                    <div className="p-4 bg-blue-50 dark:bg-blue-900/20 space-y-4">
                      <div className="flex items-center justify-between mb-2">
                        <h4 className="font-medium">Edit Submission</h4>
                      </div>

                      <div className="space-y-4">
                        <div className="space-y-2">
                          <Label>Text Submission</Label>
                          <Textarea
                            placeholder="Enter your submission text here..."
                            value={submissionText}
                            onChange={(e) => setSubmissionText(e.target.value)}
                            rows={6}
                            disabled={updateMutation.isPending || isUploading}
                          />
                        </div>

                        <div className="space-y-2">
                          <Label>File Upload</Label>
                          <div className="border-2 border-dashed rounded-lg p-6 text-center">
                            <input
                              type="file"
                              ref={editFileInputRef}
                              id="edit-file-upload"
                              className="hidden"
                              onChange={handleEditFileChange}
                              disabled={updateMutation.isPending || isUploading}
                            />
                            <label
                              htmlFor="edit-file-upload"
                              className="cursor-pointer flex flex-col items-center gap-2"
                            >
                              <Upload className="h-8 w-8 text-muted-foreground" />
                              {selectedFile ? (
                                <div>
                                  <span className="text-sm font-medium block">
                                    {selectedFile.name}
                                  </span>
                                  <span className="text-xs text-muted-foreground">
                                    {(selectedFile.size / 1024).toFixed(1)} KB
                                  </span>
                                </div>
                              ) : (
                                <>
                                  <span className="text-sm font-medium">
                                    Click to upload new file
                                  </span>
                                  <span className="text-xs text-muted-foreground">
                                    PDF, DOC, or ZIP (max 10MB)
                                  </span>
                                </>
                              )}
                            </label>
                          </div>
                        </div>

                        <div className="flex gap-2 justify-end">
                          <Button
                            variant="outline"
                            onClick={handleCancelEdit}
                            disabled={updateMutation.isPending || isUploading}
                          >
                            Cancel
                          </Button>
                          <Button
                            onClick={handleUpdateSubmit}
                            disabled={
                              updateMutation.isPending ||
                              isUploading ||
                              (!submissionText && !selectedFile)
                            }
                          >
                            {(updateMutation.isPending || isUploading) && (
                              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            )}
                            {isUploading
                              ? "Uploading..."
                              : updateMutation.isPending
                                ? "Updating..."
                                : "Save Changes"}
                          </Button>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Submission comments row */}
                  <div className="grid grid-cols-3 gap-4 p-4 bg-muted/30">
                    <div className="font-medium">Submission comments</div>
                    <div className="col-span-2">
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <MessageSquare className="h-4 w-4" />
                        <span>Comments (0)</span>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Feedback Section - Moodle Style */}
          {isStudent() && mySubmission && mySubmission.status === "graded" && (
            <Card>
              <CardHeader>
                <CardTitle>Feedback</CardTitle>
              </CardHeader>
              <CardContent className="p-0">
                <div className="divide-y">
                  {/* Grade row */}
                  <div className="grid grid-cols-3 gap-4 p-4">
                    <div className="font-medium">Grade</div>
                    <div className="col-span-2">
                      <span className="text-xl font-bold text-primary">
                        {mySubmission.grade_score?.toFixed(0)}%
                      </span>
                    </div>
                  </div>

                  {/* Graded on row */}
                  {mySubmission.graded_at && (
                    <div className="grid grid-cols-3 gap-4 p-4">
                      <div className="font-medium">Graded on</div>
                      <div className="col-span-2">
                        {formatDateTime(mySubmission.graded_at)}
                      </div>
                    </div>
                  )}

                  {/* Graded by row */}
                  {mySubmission.grader_first_name && (
                    <div className="grid grid-cols-3 gap-4 p-4">
                      <div className="font-medium">Graded by</div>
                      <div className="col-span-2 flex items-center gap-2">
                        <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center text-sm font-medium">
                          {mySubmission.grader_first_name?.charAt(0)}
                          {mySubmission.grader_last_name?.charAt(0)}
                        </div>
                        <span>
                          {getFullName(
                            mySubmission.grader_first_name,
                            mySubmission.grader_last_name,
                          )}
                        </span>
                      </div>
                    </div>
                  )}

                  {/* Feedback comments row */}
                  {mySubmission.grade_feedback && (
                    <div className="grid grid-cols-3 gap-4 p-4">
                      <div className="font-medium">Feedback comments</div>
                      <div className="col-span-2">
                        <p className="text-sm bg-muted p-3 rounded whitespace-pre-wrap">
                          {mySubmission.grade_feedback}
                        </p>
                      </div>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Teacher: Submissions List */}
          {isTeacher() && (
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Submissions</CardTitle>
                    <CardDescription>
                      {submissionsData?.total || 0} submitted,{" "}
                      {submissionsData?.submissions?.filter(
                        (s: any) => s.status === "graded",
                      ).length || 0}{" "}
                      graded
                      {selectedGroupId && courseGroups.length > 0 && (
                        <span className="ml-1">
                          in{" "}
                          {courseGroups.find(
                            (g) => g.group_id === selectedGroupId,
                          )?.group_code || "selected group"}
                        </span>
                      )}
                    </CardDescription>
                  </div>
                  {/* Group Selector for Teachers */}
                  {courseGroups.length > 0 && (
                    <div className="flex items-center gap-2">
                      <UsersRound className="h-4 w-4 text-muted-foreground" />
                      <Select
                        value={selectedGroupId || undefined}
                        onValueChange={setSelectedGroupId}
                      >
                        <SelectTrigger className="w-[140px] h-9">
                          <SelectValue placeholder="Select Group" />
                        </SelectTrigger>
                        <SelectContent>
                          {courseGroups.map((assignment) => (
                            <SelectItem
                              key={assignment.group_id}
                              value={assignment.group_id}
                            >
                              {assignment.group_code}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Badge variant="outline" className="ml-1">
                        {courseGroups.length}{" "}
                        {courseGroups.length === 1 ? "group" : "groups"}
                      </Badge>
                    </div>
                  )}
                </div>
              </CardHeader>
              <CardContent>
                {submissionsLoading ? (
                  <div className="space-y-4">
                    {[1, 2, 3].map((i) => (
                      <Skeleton key={i} className="h-16 w-full" />
                    ))}
                  </div>
                ) : submissionsData?.submissions.length ? (
                  <div className="space-y-4">
                    {submissionsData.submissions.map((submission) => (
                      <div
                        key={submission.id}
                        className="p-4 rounded-lg border bg-card"
                      >
                        {/* Student info row */}
                        <div className="flex items-center justify-between mb-3">
                          <div className="flex items-center gap-3">
                            <div className="h-10 w-10 rounded-full bg-muted flex items-center justify-center">
                              <User className="h-5 w-5 text-muted-foreground" />
                            </div>
                            <div>
                              <p className="font-medium">
                                {getFullName(
                                  submission.student_first_name,
                                  submission.student_last_name,
                                )}
                              </p>
                              <p className="text-sm text-muted-foreground">
                                {submission.student_email}
                              </p>
                            </div>
                          </div>
                          <Badge
                            className={getSubmissionStatusColor(
                              submission.status,
                            )}
                          >
                            {submission.status === "graded"
                              ? `Graded: ${submission.grade_score?.toFixed(0)}%`
                              : submission.status}
                          </Badge>
                        </div>

                        {/* Submission details */}
                        <div className="space-y-2 mb-3">
                          <p className="text-sm text-muted-foreground">
                            Submitted: {formatDateTime(submission.submitted_at)}
                            {assignment.due_at &&
                              new Date(submission.submitted_at) >
                                new Date(assignment.due_at) && (
                                <span className="text-red-500 ml-2">
                                  (Late)
                                </span>
                              )}
                          </p>

                          {/* File download */}
                          {submission.file_url && (
                            <div className="flex items-center gap-2 p-2 bg-muted/50 rounded">
                              <FileText className="h-4 w-4 text-blue-500" />
                              <span className="text-sm flex-1 truncate">
                                {submission.file_name || "Submission file"}
                              </span>
                              <Button variant="outline" size="sm" asChild>
                                <a
                                  href={submission.file_url}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  download
                                >
                                  <Download className="h-4 w-4 mr-1" />
                                  Download
                                </a>
                              </Button>
                            </div>
                          )}

                          {/* Text content preview */}
                          {submission.content_text && (
                            <div className="p-2 bg-muted/50 rounded">
                              <p className="text-sm text-muted-foreground mb-1">
                                Text submission:
                              </p>
                              <p className="text-sm whitespace-pre-wrap line-clamp-3">
                                {submission.content_text}
                              </p>
                            </div>
                          )}
                        </div>

                        {/* Grade info if graded */}
                        {submission.status === "graded" && (
                          <div className="p-3 bg-green-50 dark:bg-green-900/20 rounded mb-3">
                            <div className="flex items-center justify-between mb-2">
                              <span className="text-sm font-medium text-green-700 dark:text-green-400">
                                Grade: {submission.grade_score?.toFixed(0)}%
                              </span>
                              {submission.graded_at && (
                                <span className="text-xs text-muted-foreground">
                                  Graded: {formatDateTime(submission.graded_at)}
                                </span>
                              )}
                            </div>
                            {submission.grade_feedback && (
                              <p className="text-sm text-muted-foreground">
                                <strong>Feedback:</strong>{" "}
                                {submission.grade_feedback}
                              </p>
                            )}
                          </div>
                        )}

                        {/* Grade button */}
                        <div className="flex justify-end">
                          <Button
                            size="sm"
                            onClick={() => openGradingDialog(submission)}
                          >
                            <Star className="h-4 w-4 mr-1" />
                            {submission.status === "graded"
                              ? "Edit Grade"
                              : "Grade"}
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-center text-muted-foreground py-8">
                    No submissions yet
                  </p>
                )}
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                {mySubmission ? (
                  mySubmission.status === "graded" ? (
                    <Award className="h-5 w-5 text-green-600" />
                  ) : (
                    <CheckCircle2 className="h-5 w-5 text-blue-600" />
                  )
                ) : (
                  <Clock className="h-5 w-5 text-muted-foreground" />
                )}
                Submission Status
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {isStudent() && (
                <div className="space-y-3">
                  {/* Submission Status */}
                  <div className="flex items-center gap-3">
                    {mySubmission ? (
                      mySubmission.status === "graded" ? (
                        <>
                          <Award className="h-5 w-5 text-green-600" />
                          <div>
                            <span className="font-medium text-green-600">
                              Graded
                            </span>
                            <p className="text-sm text-muted-foreground">
                              {mySubmission.grade_score?.toFixed(0)}%
                            </p>
                          </div>
                        </>
                      ) : (
                        <>
                          <CheckCircle2 className="h-5 w-5 text-blue-600" />
                          <div>
                            <span className="font-medium text-blue-600">
                              Submitted
                            </span>
                            <p className="text-sm text-muted-foreground">
                              Awaiting grade
                            </p>
                          </div>
                        </>
                      )
                    ) : overdue ? (
                      <>
                        <XCircle className="h-5 w-5 text-destructive" />
                        <div>
                          <span className="font-medium text-destructive">
                            Not submitted
                          </span>
                          <p className="text-sm text-destructive">
                            Deadline passed
                          </p>
                        </div>
                      </>
                    ) : (
                      <>
                        <Clock className="h-5 w-5 text-yellow-600" />
                        <div>
                          <span className="font-medium text-yellow-600">
                            Not submitted
                          </span>
                          <p className="text-sm text-muted-foreground">
                            Pending
                          </p>
                        </div>
                      </>
                    )}
                  </div>
                </div>
              )}
              {isTeacher() && (
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">
                      Submissions
                    </span>
                    <span className="font-medium">
                      {submissionsData?.total || 0}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">
                      Graded
                    </span>
                    <span className="font-medium">
                      {submissionsData?.submissions.filter(
                        (s) => s.status === "graded",
                      ).length || 0}
                    </span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Grade Card - Only show if graded */}
          {isStudent() && mySubmission && mySubmission.status === "graded" && (
            <Card className="border-green-200 dark:border-green-800 bg-green-50/50 dark:bg-green-900/10">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-green-700 dark:text-green-400">
                  <Award className="h-5 w-5" />
                  Your Grade
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="text-center">
                  <p className="text-4xl font-bold text-green-600 dark:text-green-400">
                    {mySubmission.grade_score?.toFixed(0)}%
                  </p>
                  <p className="text-sm text-muted-foreground mt-2">
                    {(mySubmission.grade_score || 0) >= 90
                      ? "Excellent!"
                      : (mySubmission.grade_score || 0) >= 70
                        ? "Good work!"
                        : (mySubmission.grade_score || 0) >= 50
                          ? "Satisfactory"
                          : "Needs improvement"}
                  </p>
                </div>
                {mySubmission.grade_feedback && (
                  <div className="mt-4 pt-4 border-t">
                    <p className="text-sm font-medium mb-1">Feedback:</p>
                    <p className="text-sm text-muted-foreground bg-white dark:bg-muted p-3 rounded">
                      {mySubmission.grade_feedback}
                    </p>
                  </div>
                )}
                {mySubmission.graded_at && (
                  <p className="text-xs text-muted-foreground text-center">
                    Graded on {formatDateTime(mySubmission.graded_at)}
                  </p>
                )}
              </CardContent>
            </Card>
          )}

          {/* Deadline Card with Time Remaining */}
          <Card
            className={cn(
              overdue && !mySubmission && "border-destructive bg-destructive/5",
            )}
          >
            <CardHeader>
              <CardTitle
                className={cn(
                  "flex items-center gap-2",
                  overdue && !mySubmission && "text-destructive",
                )}
              >
                <Calendar className="h-5 w-5" />
                Deadline
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <p className="text-lg font-medium">
                {assignment.due_at
                  ? formatDateTime(assignment.due_at)
                  : "No deadline"}
              </p>

              {/* Time Remaining */}
              {assignment.due_at && !mySubmission && (
                <div
                  className={cn(
                    "p-3 rounded-lg",
                    overdue
                      ? "bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300"
                      : "bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300",
                  )}
                >
                  <div className="flex items-center gap-2">
                    {overdue ? (
                      <AlertTriangle className="h-4 w-4" />
                    ) : (
                      <Clock className="h-4 w-4" />
                    )}
                    <span className="font-medium">
                      {getTimeRemaining(assignment.due_at)}
                    </span>
                  </div>
                </div>
              )}

              {/* Show submission time info if submitted */}
              {assignment.due_at && mySubmission && (
                <div
                  className={cn(
                    "p-3 rounded-lg",
                    new Date(mySubmission.submitted_at) >
                      new Date(assignment.due_at)
                      ? "bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300"
                      : "bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300",
                  )}
                >
                  <div className="flex items-center gap-2">
                    {new Date(mySubmission.submitted_at) >
                    new Date(assignment.due_at) ? (
                      <AlertTriangle className="h-4 w-4" />
                    ) : (
                      <CheckCircle2 className="h-4 w-4" />
                    )}
                    <span className="font-medium text-sm">
                      {getSubmissionTimeInfo(
                        mySubmission.submitted_at,
                        assignment.due_at,
                      )}
                    </span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Teacher Grading Dialog */}
      <Dialog
        open={!!gradingSubmission}
        onOpenChange={(open) => !open && setGradingSubmission(null)}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Star className="h-5 w-5 text-yellow-500" />
              Grade Submission
            </DialogTitle>
            <DialogDescription>
              {gradingSubmission && (
                <>
                  Student:{" "}
                  <strong>
                    {getFullName(
                      gradingSubmission.student_first_name,
                      gradingSubmission.student_last_name,
                    )}
                  </strong>
                </>
              )}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            {/* Download file if exists */}
            {gradingSubmission?.file_url && (
              <div className="flex items-center gap-2 p-3 bg-muted rounded-lg">
                <FileText className="h-5 w-5 text-blue-500" />
                <span className="flex-1 text-sm truncate">
                  {gradingSubmission.file_name || "Submission file"}
                </span>
                <Button variant="outline" size="sm" asChild>
                  <a
                    href={gradingSubmission.file_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    download
                  >
                    <Download className="h-4 w-4" />
                  </a>
                </Button>
              </div>
            )}

            {/* Text content if exists */}
            {gradingSubmission?.content_text && (
              <div className="p-3 bg-muted rounded-lg max-h-32 overflow-y-auto">
                <p className="text-xs text-muted-foreground mb-1">
                  Text submission:
                </p>
                <p className="text-sm whitespace-pre-wrap">
                  {gradingSubmission.content_text}
                </p>
              </div>
            )}

            {/* Score input */}
            <div className="space-y-2">
              <Label htmlFor="grade-score">Score (percentage)</Label>
              <div className="flex items-center gap-2">
                <Input
                  id="grade-score"
                  type="number"
                  min="0"
                  max={100}
                  step="1"
                  value={gradeScore}
                  onChange={(e) => setGradeScore(e.target.value)}
                  placeholder="Enter percentage..."
                  className="flex-1"
                />
                <span className="text-muted-foreground">%</span>
              </div>
            </div>

            {/* Feedback textarea */}
            <div className="space-y-2">
              <Label htmlFor="grade-feedback">Feedback (optional)</Label>
              <Textarea
                id="grade-feedback"
                placeholder="Enter feedback for the student..."
                value={gradeFeedback}
                onChange={(e) => setGradeFeedback(e.target.value)}
                rows={4}
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setGradingSubmission(null)}
              disabled={gradeMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              onClick={handleGradeSubmit}
              disabled={gradeMutation.isPending || !gradeScore}
            >
              {gradeMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              {gradingSubmission?.status === "graded"
                ? "Update Grade"
                : "Submit Grade"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </motion.div>
  );
}
