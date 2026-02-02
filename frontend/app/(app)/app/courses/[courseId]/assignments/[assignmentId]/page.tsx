"use client";

import { useState } from "react";
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

export default function AssignmentDetailPage() {
  const params = useParams();
  const courseId = params.courseId as string;
  const assignmentId = params.assignmentId as string;
  const queryClient = useQueryClient();
  const { isTeacher, isStudent } = useAuthStore();

  const [submissionText, setSubmissionText] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const { data: assignment, isLoading: assignmentLoading } = useQuery({
    queryKey: ["assignment", assignmentId],
    queryFn: () => assignmentService.get(assignmentId),
  });

  const { data: submissionsData, isLoading: submissionsLoading } = useQuery({
    queryKey: ["assignment-submissions", assignmentId],
    queryFn: () => submissionService.getByAssignment(assignmentId),
    enabled: isTeacher(),
  });

  const { data: mySubmissions } = useQuery({
    queryKey: ["my-submissions"],
    queryFn: () => submissionService.getMySubmissions(),
    enabled: isStudent(),
  });

  const mySubmission = mySubmissions?.submissions.find(
    (s) => s.assignment_id === assignmentId,
  );

  const submitMutation = useMutation({
    mutationFn: () =>
      submissionService.create({
        assignment_id: assignmentId,
        content_text: submissionText,
        file_url: selectedFile ? URL.createObjectURL(selectedFile) : undefined,
      }),
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

  const handleSubmit = () => {
    if (!submissionText && !selectedFile) {
      toast.error("Please provide either text or a file");
      return;
    }
    submitMutation.mutate();
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
                  {assignment.description || "No description provided."}
                </p>
              </div>
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
                  <span>{assignment.max_points} points</span>
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
                    (!submissionText && !selectedFile)
                  }
                  className="w-full"
                >
                  {submitMutation.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Submit Assignment
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
                  <div className={cn(
                    "grid grid-cols-3 gap-4 p-4",
                    mySubmission.status === "graded" 
                      ? "bg-green-50 dark:bg-green-900/20" 
                      : "bg-yellow-50 dark:bg-yellow-900/20"
                  )}>
                    <div className="font-medium">Submission status</div>
                    <div className="col-span-2">
                      {mySubmission.status === "graded" 
                        ? "Submitted for grading" 
                        : "Submitted for grading"}
                    </div>
                  </div>
                  
                  {/* Grading status row */}
                  <div className={cn(
                    "grid grid-cols-3 gap-4 p-4",
                    mySubmission.status === "graded"
                      ? "bg-green-50 dark:bg-green-900/20"
                      : "bg-yellow-50 dark:bg-yellow-900/20"
                  )}>
                    <div className="font-medium">Grading status</div>
                    <div className="col-span-2">
                      {mySubmission.status === "graded" ? (
                        <span className="text-green-600 dark:text-green-400 font-medium">Graded</span>
                      ) : (
                        <span className="text-yellow-600 dark:text-yellow-400">Not graded</span>
                      )}
                    </div>
                  </div>
                  
                  {/* Time remaining row */}
                  <div className={cn(
                    "grid grid-cols-3 gap-4 p-4",
                    assignment.due_at && new Date(mySubmission.submitted_at) > new Date(assignment.due_at)
                      ? "bg-red-50 dark:bg-red-900/20"
                      : "bg-green-50 dark:bg-green-900/20"
                  )}>
                    <div className="font-medium">Time remaining</div>
                    <div className={cn(
                      "col-span-2",
                      assignment.due_at && new Date(mySubmission.submitted_at) > new Date(assignment.due_at)
                        ? "text-red-600 dark:text-red-400"
                        : "text-green-600 dark:text-green-400"
                    )}>
                      {getSubmissionTimeInfo(mySubmission.submitted_at, assignment.due_at)}
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
                        <a 
                          href={mySubmission.file_url} 
                          target="_blank" 
                          rel="noopener noreferrer"
                          className="flex items-center gap-2 text-primary hover:underline"
                        >
                          <FileText className="h-4 w-4" />
                          <span>Download submission</span>
                          <Download className="h-3 w-3" />
                        </a>
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
                        {mySubmission.grade_score?.toFixed(2)} / {assignment.max_points}.00
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
                          {mySubmission.grader_first_name?.charAt(0)}{mySubmission.grader_last_name?.charAt(0)}
                        </div>
                        <span>
                          {getFullName(mySubmission.grader_first_name, mySubmission.grader_last_name)}
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
                <CardTitle>Submissions</CardTitle>
                <CardDescription>
                  {submissionsData?.total || 0} students have submitted
                </CardDescription>
              </CardHeader>
              <CardContent>
                {submissionsLoading ? (
                  <div className="space-y-4">
                    {[1, 2, 3].map((i) => (
                      <Skeleton key={i} className="h-16 w-full" />
                    ))}
                  </div>
                ) : submissionsData?.submissions.length ? (
                  <div className="space-y-3">
                    {submissionsData.submissions.map((submission) => (
                      <div
                        key={submission.id}
                        className="flex items-center justify-between p-4 rounded-lg border"
                      >
                        <div className="flex items-center gap-3">
                          <div className="h-10 w-10 rounded-full bg-muted flex items-center justify-center">
                            {submission.status === "graded" ? (
                              <CheckCircle2 className="h-5 w-5 text-green-600" />
                            ) : (
                              <FileText className="h-5 w-5 text-muted-foreground" />
                            )}
                          </div>
                          <div>
                            <p className="font-medium">
                              {getFullName(
                                submission.student_first_name,
                                submission.student_last_name,
                              )}
                            </p>
                            <p className="text-sm text-muted-foreground">
                              Submitted{" "}
                              {formatDateTime(submission.submitted_at)}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge
                            className={getSubmissionStatusColor(
                              submission.status,
                            )}
                          >
                            {submission.status}
                          </Badge>
                          {submission.grade_score !== undefined && (
                            <span className="font-medium">
                              {submission.grade_score}/{assignment.max_points}
                            </span>
                          )}
                          <Button size="sm" variant="outline">
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
                            <span className="font-medium text-green-600">Graded</span>
                            <p className="text-sm text-muted-foreground">
                              {mySubmission.grade_score}/{assignment.max_points} points
                            </p>
                          </div>
                        </>
                      ) : (
                        <>
                          <CheckCircle2 className="h-5 w-5 text-blue-600" />
                          <div>
                            <span className="font-medium text-blue-600">Submitted</span>
                            <p className="text-sm text-muted-foreground">Awaiting grade</p>
                          </div>
                        </>
                      )
                    ) : overdue ? (
                      <>
                        <XCircle className="h-5 w-5 text-destructive" />
                        <div>
                          <span className="font-medium text-destructive">Not submitted</span>
                          <p className="text-sm text-destructive">Deadline passed</p>
                        </div>
                      </>
                    ) : (
                      <>
                        <Clock className="h-5 w-5 text-yellow-600" />
                        <div>
                          <span className="font-medium text-yellow-600">Not submitted</span>
                          <p className="text-sm text-muted-foreground">Pending</p>
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
                    {mySubmission.grade_score}
                  </p>
                  <p className="text-lg text-muted-foreground">
                    out of {assignment.max_points}
                  </p>
                  <p className="text-sm text-muted-foreground mt-2">
                    {Math.round((mySubmission.grade_score || 0) / assignment.max_points * 100)}%
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
          <Card className={cn(
            overdue && !mySubmission && "border-destructive bg-destructive/5"
          )}>
            <CardHeader>
              <CardTitle className={cn(
                "flex items-center gap-2",
                overdue && !mySubmission && "text-destructive"
              )}>
                <Calendar className="h-5 w-5" />
                Deadline
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <p className="text-lg font-medium">
                {assignment.due_at ? formatDateTime(assignment.due_at) : "No deadline"}
              </p>
              
              {/* Time Remaining */}
              {assignment.due_at && !mySubmission && (
                <div className={cn(
                  "p-3 rounded-lg",
                  overdue 
                    ? "bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300"
                    : "bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300"
                )}>
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
                <div className={cn(
                  "p-3 rounded-lg",
                  new Date(mySubmission.submitted_at) > new Date(assignment.due_at)
                    ? "bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300"
                    : "bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300"
                )}>
                  <div className="flex items-center gap-2">
                    {new Date(mySubmission.submitted_at) > new Date(assignment.due_at) ? (
                      <AlertTriangle className="h-4 w-4" />
                    ) : (
                      <CheckCircle2 className="h-4 w-4" />
                    )}
                    <span className="font-medium text-sm">
                      {getSubmissionTimeInfo(mySubmission.submitted_at, assignment.due_at)}
                    </span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </motion.div>
  );
}
