"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import {
  Calendar,
  Plus,
  Trash2,
  Users,
  Check,
  X,
  Clock,
  AlertCircle,
  Loader2,
  ChevronRight,
} from "lucide-react";
import { toast } from "sonner";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Progress } from "@/components/ui/progress";

import { attendanceService } from "@/services/attendance";
import groupService from "@/services/groups";
import type { AttendanceStatus, AttendanceSession } from "@/types/attendance";
import api from "@/lib/api-client";

// ==================== STUDENT VIEW ====================
export function StudentAttendanceView({ courseId }: { courseId: string }) {
  const { data, isLoading } = useQuery({
    queryKey: ["my-attendance", courseId],
    queryFn: () => attendanceService.getMyAttendance(courseId),
  });

  const summary = data?.summary;
  const marks = data?.marks || [];

  if (isLoading) {
    return (
      <Card>
        <CardContent className="p-8 text-center text-muted-foreground">
          Loading attendance...
        </CardContent>
      </Card>
    );
  }

  const statusConfig: Record<
    string,
    { label: string; color: string; icon: React.ElementType }
  > = {
    present: {
      label: "Present",
      color:
        "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
      icon: Check,
    },
    absent: {
      label: "Absent",
      color: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
      icon: X,
    },
    late: {
      label: "Late",
      color:
        "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
      icon: Clock,
    },
    excused: {
      label: "Excused",
      color: "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
      icon: AlertCircle,
    },
  };

  return (
    <div className="space-y-4">
      {/* Summary Card */}
      <Card>
        <CardHeader>
          <CardTitle>Attendance</CardTitle>
          <CardDescription>
            Your attendance record for this course
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {/* Attendance Rate */}
            <div className="flex justify-between items-center p-4 rounded-lg bg-muted">
              <div>
                <span className="font-medium">Attendance Rate</span>
                <p className="text-sm text-muted-foreground">
                  {summary?.total_classes || 0} classes total
                </p>
              </div>
              <span
                className={`text-3xl font-bold ${
                  (summary?.percentage || 0) >= 80
                    ? "text-green-600"
                    : (summary?.percentage || 0) >= 60
                      ? "text-yellow-600"
                      : "text-red-600"
                }`}
              >
                {(summary?.percentage || 0).toFixed(1)}%
              </span>
            </div>

            {/* Progress bar */}
            <Progress value={summary?.percentage || 0} className="h-2" />

            {/* Stats row */}
            {summary && summary.total_classes > 0 && (
              <div className="grid grid-cols-4 gap-2 text-center">
                <div className="p-2 rounded bg-green-50 dark:bg-green-900/20">
                  <p className="text-lg font-bold text-green-600">
                    {summary.present}
                  </p>
                  <p className="text-xs text-muted-foreground">Present</p>
                </div>
                <div className="p-2 rounded bg-red-50 dark:bg-red-900/20">
                  <p className="text-lg font-bold text-red-600">
                    {summary.absent}
                  </p>
                  <p className="text-xs text-muted-foreground">Absent</p>
                </div>
                <div className="p-2 rounded bg-yellow-50 dark:bg-yellow-900/20">
                  <p className="text-lg font-bold text-yellow-600">
                    {summary.late}
                  </p>
                  <p className="text-xs text-muted-foreground">Late</p>
                </div>
                <div className="p-2 rounded bg-blue-50 dark:bg-blue-900/20">
                  <p className="text-lg font-bold text-blue-600">
                    {summary.excused}
                  </p>
                  <p className="text-xs text-muted-foreground">Excused</p>
                </div>
              </div>
            )}

            {/* Records list */}
            {marks.length > 0 ? (
              <div className="space-y-2">
                <h4 className="text-sm font-medium text-muted-foreground">
                  History
                </h4>
                {marks.map((mark) => {
                  const cfg = statusConfig[mark.status] || statusConfig.absent;
                  const Icon = cfg.icon;
                  return (
                    <div
                      key={mark.id}
                      className="flex justify-between items-center p-3 rounded border"
                    >
                      <span className="text-sm">
                        {format(new Date(mark.marked_at), "MMM d, yyyy")}
                      </span>
                      <Badge className={cfg.color}>
                        <Icon className="h-3 w-3 mr-1" />
                        {cfg.label}
                      </Badge>
                    </div>
                  );
                })}
              </div>
            ) : (
              <p className="text-center text-muted-foreground py-4">
                No attendance records yet
              </p>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

// ==================== TEACHER VIEW ====================
interface TeacherAttendanceViewProps {
  courseId: string;
  courseGroups?: { group_id: string; group_code: string }[];
}

export function TeacherAttendanceView({
  courseId,
  courseGroups: propGroups,
}: TeacherAttendanceViewProps) {
  const queryClient = useQueryClient();
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [selectedDate, setSelectedDate] = useState(
    format(new Date(), "yyyy-MM-dd"),
  );
  const [startTime, setStartTime] = useState("16:00");
  const [endTime, setEndTime] = useState("16:50");
  const [selectedGroupForAttendance, setSelectedGroupForAttendance] =
    useState<string>("");
  const [activeSession, setActiveSession] = useState<AttendanceSession | null>(
    null,
  );
  const [markingAttendance, setMarkingAttendance] = useState<
    Record<string, AttendanceStatus>
  >({});

  // Fetch groups for this course directly
  const { data: courseGroupsData, isLoading: groupsLoading } = useQuery({
    queryKey: ["course-groups", courseId],
    queryFn: () => api.get<any>(`/groups/course/${courseId}`),
  });

  // Also fetch enrolled students as fallback
  const { data: enrolledStudents } = useQuery({
    queryKey: ["enrolled-students", courseId],
    queryFn: () => api.get<any>(`/enrollments/course/${courseId}?take=200`),
  });

  // Use fetched groups, mapping to expected shape
  const rawGroups = Array.isArray(courseGroupsData)
    ? courseGroupsData
    : courseGroupsData?.groups || courseGroupsData?.data || [];
  const courseGroups = (Array.isArray(rawGroups) ? rawGroups : []).map(
    (g: any) => ({
      group_id: g.id,
      group_code: g.name || g.id?.slice(0, 8) || "Unknown",
    }),
  );

  const hasGroups = courseGroups.length > 0;

  // Fetch sessions for this course
  const { data: sessionsData, isLoading: sessionsLoading } = useQuery({
    queryKey: ["attendance-sessions", courseId],
    queryFn: () => attendanceService.listSessions(courseId),
  });

  // Fetch group members when a group is selected for taking attendance
  const { data: groupMembersData } = useQuery({
    queryKey: ["group-members", selectedGroupForAttendance],
    queryFn: () => groupService.getMembers(selectedGroupForAttendance),
    enabled: !!selectedGroupForAttendance && hasGroups,
  });

  // Build a unified student list: group members OR enrolled students
  const groupMembers = (() => {
    if (hasGroups && selectedGroupForAttendance && groupMembersData) {
      const raw = Array.isArray(groupMembersData)
        ? groupMembersData
        : (groupMembersData as any)?.members || [];
      return Array.isArray(raw) ? raw : [];
    }
    if (!hasGroups && enrolledStudents) {
      const raw = Array.isArray(enrolledStudents)
        ? enrolledStudents
        : (enrolledStudents as any)?.enrollments ||
          (enrolledStudents as any)?.data ||
          [];
      return (Array.isArray(raw) ? raw : []).map((e: any) => ({
        student_id: e.student_id || e.user_id || e.id,
        student_first_name: e.first_name || "",
        student_last_name: e.last_name || "",
        student_email: e.email || "",
      }));
    }
    return [];
  })();

  console.log(
    "[ATTENDANCE] courseGroups:",
    courseGroups.length,
    "hasGroups:",
    hasGroups,
    "enrolledStudents:",
    enrolledStudents,
    "groupMembers:",
    groupMembers.length,
  );

  // Fetch marks when viewing a session
  const { data: sessionMarks } = useQuery({
    queryKey: ["session-marks", activeSession?.id],
    queryFn: () => attendanceService.getMarksBySession(activeSession!.id),
    enabled: !!activeSession,
  });

  // Create session mutation
  const createSessionMutation = useMutation({
    mutationFn: (data: {
      course_id: string;
      date: string;
      start_time?: string;
      end_time?: string;
    }) => attendanceService.createSession(data),
    onSuccess: (session) => {
      toast.success("Attendance session created");
      queryClient.invalidateQueries({
        queryKey: ["attendance-sessions", courseId],
      });
      setShowCreateDialog(false);
      setActiveSession(session);
    },
    onError: (err: Error) =>
      toast.error(err.message || "Failed to create session"),
  });

  // Bulk mark mutation
  const bulkMarkMutation = useMutation({
    mutationFn: (data: {
      session_id: string;
      marks: { student_id: string; status: AttendanceStatus }[];
    }) => attendanceService.bulkMark(data),
    onSuccess: () => {
      toast.success("Attendance saved!");
      queryClient.invalidateQueries({
        queryKey: ["attendance-sessions", courseId],
      });
      queryClient.invalidateQueries({
        queryKey: ["session-marks", activeSession?.id],
      });
      setActiveSession(null);
      setMarkingAttendance({});
    },
    onError: (err: Error) =>
      toast.error(err.message || "Failed to save attendance"),
  });

  // Delete session mutation
  const deleteSessionMutation = useMutation({
    mutationFn: (id: string) => attendanceService.deleteSession(id),
    onSuccess: () => {
      toast.success("Session deleted");
      queryClient.invalidateQueries({
        queryKey: ["attendance-sessions", courseId],
      });
    },
    onError: (err: Error) => toast.error(err.message || "Failed to delete"),
  });

  const handleCreateAndMark = () => {
    if (!selectedDate) {
      toast.error("Please select a date");
      return;
    }
    createSessionMutation.mutate({
      course_id: courseId,
      date: selectedDate,
      start_time: startTime || undefined,
      end_time: endTime || undefined,
    });
  };

  const handleSaveAttendance = () => {
    if (!activeSession) return;
    const marks = Object.entries(markingAttendance).map(
      ([student_id, status]) => ({
        student_id,
        status,
      }),
    );
    if (marks.length === 0) {
      toast.error("Please mark at least one student");
      return;
    }
    bulkMarkMutation.mutate({ session_id: activeSession.id, marks });
  };

  const setAllStatus = (status: AttendanceStatus) => {
    if (!groupMembers) return;
    const newMarks: Record<string, AttendanceStatus> = {};
    groupMembers.forEach((m: any) => {
      newMarks[m.student_id] = status;
    });
    setMarkingAttendance(newMarks);
  };

  const sessions = sessionsData?.sessions || [];

  // If teacher is marking attendance for a session
  if (activeSession) {
    const members = groupMembers || [];
    const existingMarks = sessionMarks?.marks || [];

    // Pre-fill from existing marks
    if (
      existingMarks.length > 0 &&
      Object.keys(markingAttendance).length === 0
    ) {
      const prefilled: Record<string, AttendanceStatus> = {};
      existingMarks.forEach((m) => {
        prefilled[m.student_id] = m.status;
      });
      // Use members list to show all students, fill existing statuses
      if (members.length > 0) {
        members.forEach((m: any) => {
          if (!prefilled[m.student_id]) {
            prefilled[m.student_id] = "present";
          }
        });
      }
      if (Object.keys(prefilled).length > 0) {
        setMarkingAttendance(prefilled);
        return null; // Re-render with prefilled data
      }
    }

    return (
      <div className="space-y-4">
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <Calendar className="h-5 w-5" />
                  Mark Attendance -{" "}
                  {format(new Date(activeSession.starts_at), "MMMM d, yyyy")}
                </CardTitle>
                <CardDescription>
                  {activeSession.ends_at
                    ? `${format(new Date(activeSession.starts_at), "HH:mm")} — ${format(new Date(activeSession.ends_at), "HH:mm")}`
                    : format(new Date(activeSession.starts_at), "HH:mm") !==
                        "00:00"
                      ? `Starts at ${format(new Date(activeSession.starts_at), "HH:mm")}`
                      : "Select a group and mark each student's attendance"}
                </CardDescription>
              </div>
              <Button
                variant="outline"
                onClick={() => {
                  setActiveSession(null);
                  setMarkingAttendance({});
                }}
              >
                Back
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {/* Group selector - only show if groups exist */}
              {hasGroups ? (
                <div className="flex items-center gap-4">
                  <div className="flex-1">
                    <Label>Select Group</Label>
                    <Select
                      value={selectedGroupForAttendance}
                      onValueChange={(v) => {
                        setSelectedGroupForAttendance(v);
                        setMarkingAttendance({});
                      }}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="Choose a group..." />
                      </SelectTrigger>
                      <SelectContent position="popper" className="z-[100]">
                        {courseGroups.map((g) => (
                          <SelectItem key={g.group_id} value={g.group_id}>
                            {g.group_code}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>

                  {/* Quick set all */}
                  {members.length > 0 && (
                    <div className="flex gap-1 pt-5">
                      <Button
                        size="sm"
                        variant="outline"
                        className="text-green-600"
                        onClick={() => setAllStatus("present")}
                      >
                        All Present
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        className="text-red-600"
                        onClick={() => setAllStatus("absent")}
                      >
                        All Absent
                      </Button>
                    </div>
                  )}
                </div>
              ) : (
                <div className="flex items-center justify-between">
                  <p className="text-sm text-muted-foreground">
                    Showing all enrolled students ({members.length})
                  </p>
                  {members.length > 0 && (
                    <div className="flex gap-1">
                      <Button
                        size="sm"
                        variant="outline"
                        className="text-green-600"
                        onClick={() => setAllStatus("present")}
                      >
                        All Present
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        className="text-red-600"
                        onClick={() => setAllStatus("absent")}
                      >
                        All Absent
                      </Button>
                    </div>
                  )}
                </div>
              )}

              {/* Students list */}
              {hasGroups && !selectedGroupForAttendance ? (
                <p className="text-center text-muted-foreground py-8">
                  Please select a group to see students
                </p>
              ) : members.length === 0 ? (
                <p className="text-center text-muted-foreground py-8">
                  No students in this group
                </p>
              ) : (
                <div className="space-y-2">
                  {members.map((member: any, idx: number) => {
                    const studentId = member.student_id;
                    const currentStatus =
                      markingAttendance[studentId] || "present";
                    return (
                      <div
                        key={studentId}
                        className="flex items-center justify-between p-3 rounded-lg border"
                      >
                        <div className="flex items-center gap-3">
                          <span className="text-sm text-muted-foreground w-6">
                            {idx + 1}
                          </span>
                          <div>
                            <p className="font-medium text-sm">
                              {member.student_first_name ||
                                member.student_email ||
                                studentId.slice(0, 8)}
                              {member.student_last_name
                                ? ` ${member.student_last_name}`
                                : ""}
                            </p>
                            {member.student_email && (
                              <p className="text-xs text-muted-foreground">
                                {member.student_email}
                              </p>
                            )}
                          </div>
                        </div>

                        <div className="flex gap-1">
                          {(
                            [
                              "present",
                              "late",
                              "absent",
                              "excused",
                            ] as AttendanceStatus[]
                          ).map((status) => {
                            const isActive = currentStatus === status;
                            const colorMap: Record<string, string> = {
                              present: isActive
                                ? "bg-green-600 text-white"
                                : "hover:bg-green-100",
                              absent: isActive
                                ? "bg-red-600 text-white"
                                : "hover:bg-red-100",
                              late: isActive
                                ? "bg-yellow-600 text-white"
                                : "hover:bg-yellow-100",
                              excused: isActive
                                ? "bg-blue-600 text-white"
                                : "hover:bg-blue-100",
                            };
                            return (
                              <Button
                                key={status}
                                size="sm"
                                variant={isActive ? "default" : "outline"}
                                className={`text-xs px-2 h-7 ${colorMap[status]}`}
                                onClick={() =>
                                  setMarkingAttendance((prev) => ({
                                    ...prev,
                                    [studentId]: status,
                                  }))
                                }
                              >
                                {status.charAt(0).toUpperCase() +
                                  status.slice(1)}
                              </Button>
                            );
                          })}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}

              {/* Save button */}
              {members.length > 0 &&
                (!hasGroups || selectedGroupForAttendance) && (
                  <Button
                    className="w-full"
                    onClick={handleSaveAttendance}
                    disabled={bulkMarkMutation.isPending}
                  >
                    {bulkMarkMutation.isPending && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Save Attendance ({Object.keys(markingAttendance).length}{" "}
                    students)
                  </Button>
                )}
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Main sessions list view
  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-5 w-5" />
                Attendance Sessions
              </CardTitle>
              <CardDescription>
                {sessions.length} session(s) recorded
              </CardDescription>
            </div>
            <Button onClick={() => setShowCreateDialog(true)}>
              <Plus className="h-4 w-4 mr-2" />
              New Session
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {sessionsLoading ? (
            <p className="text-center text-muted-foreground py-4">Loading...</p>
          ) : sessions.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">
              No attendance sessions yet. Click &ldquo;New Session&rdquo; to
              create one.
            </p>
          ) : (
            <div className="space-y-2">
              {sessions.map((session) => (
                <div
                  key={session.id}
                  className="flex items-center justify-between p-3 rounded-lg border hover:bg-muted/50 transition-colors"
                >
                  <div
                    className="flex items-center gap-3 flex-1 cursor-pointer"
                    onClick={() => setActiveSession(session)}
                  >
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="font-medium text-sm">
                        {format(
                          new Date(session.starts_at),
                          "EEEE, MMMM d, yyyy",
                        )}
                        {session.ends_at
                          ? ` · ${format(new Date(session.starts_at), "HH:mm")} – ${format(new Date(session.ends_at), "HH:mm")}`
                          : format(new Date(session.starts_at), "HH:mm") !==
                              "00:00"
                            ? ` · ${format(new Date(session.starts_at), "HH:mm")}`
                            : ""}
                      </p>
                      <div className="flex gap-2 mt-1">
                        {session.total_students > 0 && (
                          <>
                            <Badge
                              variant="outline"
                              className="text-xs bg-green-50 text-green-700"
                            >
                              {session.present_count} present
                            </Badge>
                            {session.absent_count > 0 && (
                              <Badge
                                variant="outline"
                                className="text-xs bg-red-50 text-red-700"
                              >
                                {session.absent_count} absent
                              </Badge>
                            )}
                            {session.late_count > 0 && (
                              <Badge
                                variant="outline"
                                className="text-xs bg-yellow-50 text-yellow-700"
                              >
                                {session.late_count} late
                              </Badge>
                            )}
                          </>
                        )}
                        {session.total_students === 0 && (
                          <span className="text-xs text-muted-foreground">
                            No marks yet
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="flex gap-1">
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => setActiveSession(session)}
                    >
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                    <Button
                      size="sm"
                      variant="ghost"
                      className="text-red-500"
                      onClick={() => deleteSessionMutation.mutate(session.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Create Session Dialog */}
      <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>New Attendance Session</DialogTitle>
            <DialogDescription>
              Select a date and time range for the session
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div>
              <Label htmlFor="att-date">Date</Label>
              <Input
                id="att-date"
                type="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="att-start-time">Start Time</Label>
                <Input
                  id="att-start-time"
                  type="time"
                  value={startTime}
                  onChange={(e) => setStartTime(e.target.value)}
                />
              </div>
              <div>
                <Label htmlFor="att-end-time">End Time</Label>
                <Input
                  id="att-end-time"
                  type="time"
                  value={endTime}
                  onChange={(e) => setEndTime(e.target.value)}
                />
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setShowCreateDialog(false)}
            >
              Cancel
            </Button>
            <Button
              onClick={handleCreateAndMark}
              disabled={createSessionMutation.isPending}
            >
              {createSessionMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Create & Mark Attendance
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
