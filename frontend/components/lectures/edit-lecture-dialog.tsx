"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  CalendarIcon,
  Loader2,
  Clock,
  MapPin,
  Video,
  Pencil,
  Trash2,
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
import { Switch } from "@/components/ui/switch";
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
import { scheduleService } from "@/services/schedule";
import type {
  ScheduleEvent,
  EventType,
  RecurrenceType,
} from "@/types/schedule";

interface EditLectureDialogProps {
  lecture: ScheduleEvent;
  courseId: string;
  trigger?: React.ReactNode;
}

export function EditLectureDialog({
  lecture,
  courseId,
  trigger,
}: EditLectureDialogProps) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState(lecture.title);
  const [description, setDescription] = useState(lecture.description || "");
  const [eventType, setEventType] = useState<EventType>(lecture.event_type);
  const [startDate, setStartDate] = useState<Date | undefined>(
    new Date(lecture.start_time),
  );
  const [startTime, setStartTime] = useState(
    format(new Date(lecture.start_time), "HH:mm"),
  );
  const [endTime, setEndTime] = useState(
    format(new Date(lecture.end_time), "HH:mm"),
  );
  const [location, setLocation] = useState(lecture.location || "");
  const [isOnline, setIsOnline] = useState(lecture.is_online || false);
  const [meetingUrl, setMeetingUrl] = useState(lecture.meeting_url || "");
  const [recurrence, setRecurrence] = useState<RecurrenceType>(
    lecture.recurrence || "none",
  );
  const [errors, setErrors] = useState<{ title?: string; date?: string }>({});
  const queryClient = useQueryClient();

  const handleOpenChange = (isOpen: boolean) => {
    setOpen(isOpen);
    if (isOpen) {
      // Reset form with current lecture data when opening
      setTitle(lecture.title);
      setDescription(lecture.description || "");
      setEventType(lecture.event_type);
      setStartDate(new Date(lecture.start_time));
      setStartTime(format(new Date(lecture.start_time), "HH:mm"));
      setEndTime(format(new Date(lecture.end_time), "HH:mm"));
      setLocation(lecture.location || "");
      setIsOnline(lecture.is_online || false);
      setMeetingUrl(lecture.meeting_url || "");
      setRecurrence(lecture.recurrence || "none");
      setErrors({});
    }
  };

  const updateMutation = useMutation({
    mutationFn: () => {
      if (!startDate) throw new Error("Start date is required");

      const startDateTime = new Date(startDate);
      const [startHour, startMinute] = startTime.split(":").map(Number);
      startDateTime.setHours(startHour, startMinute, 0, 0);

      const endDateTime = new Date(startDate);
      const [endHour, endMinute] = endTime.split(":").map(Number);
      endDateTime.setHours(endHour, endMinute, 0, 0);

      return scheduleService.update(lecture.id, {
        title,
        description: description || undefined,
        event_type: eventType,
        start_time: startDateTime.toISOString(),
        end_time: endDateTime.toISOString(),
        location: location || undefined,
        is_online: isOnline,
        meeting_url: isOnline && meetingUrl ? meetingUrl : undefined,
        recurrence,
      });
    },
    onSuccess: () => {
      toast.success("Lecture updated successfully!");
      queryClient.invalidateQueries({
        queryKey: ["course-schedule", courseId],
      });
      queryClient.invalidateQueries({ queryKey: ["schedule"] });
      setOpen(false);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to update lecture");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: () => scheduleService.delete(lecture.id),
    onSuccess: () => {
      toast.success("Lecture deleted successfully!");
      queryClient.invalidateQueries({
        queryKey: ["course-schedule", courseId],
      });
      queryClient.invalidateQueries({ queryKey: ["schedule"] });
      setOpen(false);
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to delete lecture");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const newErrors: { title?: string; date?: string } = {};
    if (!title || title.length < 2) {
      newErrors.title = "Title must be at least 2 characters";
    }
    if (!startDate) {
      newErrors.date = "Please select a date";
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
      <DialogContent className="sm:max-w-137.5 max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Edit Lecture</DialogTitle>
          <DialogDescription>
            Update lecture details or delete this lecture.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="edit-title">Title *</Label>
            <Input
              id="edit-title"
              placeholder="e.g., Introduction to Programming"
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                if (errors.title)
                  setErrors((prev) => ({ ...prev, title: undefined }));
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
              placeholder="Lecture description or topics to cover..."
              className="resize-none min-h-20"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={isPending}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Event Type</Label>
              <Select
                value={eventType}
                onValueChange={(v) => setEventType(v as EventType)}
                disabled={isPending}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="class">Lecture</SelectItem>
                  <SelectItem value="lab">Lab Session</SelectItem>
                  <SelectItem value="exam">Exam</SelectItem>
                  <SelectItem value="office_hours">Office Hours</SelectItem>
                  <SelectItem value="event">Other Event</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>Recurrence</Label>
              <Select
                value={recurrence}
                onValueChange={(v) => setRecurrence(v as RecurrenceType)}
                disabled={isPending}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">One-time</SelectItem>
                  <SelectItem value="daily">Daily</SelectItem>
                  <SelectItem value="weekly">Weekly</SelectItem>
                  <SelectItem value="monthly">Monthly</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label>Date *</Label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  className={cn(
                    "w-full pl-3 text-left font-normal",
                    !startDate && "text-muted-foreground",
                  )}
                  disabled={isPending}
                >
                  {startDate ? (
                    format(startDate, "PPP")
                  ) : (
                    <span>Pick a date</span>
                  )}
                  <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0" align="start">
                <Calendar
                  mode="single"
                  selected={startDate}
                  onSelect={(date) => {
                    setStartDate(date);
                    if (errors.date)
                      setErrors((prev) => ({ ...prev, date: undefined }));
                  }}
                  initialFocus
                />
              </PopoverContent>
            </Popover>
            {errors.date && (
              <p className="text-sm text-destructive">{errors.date}</p>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="edit-start-time">Start Time</Label>
              <div className="relative">
                <Clock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="edit-start-time"
                  type="time"
                  value={startTime}
                  onChange={(e) => setStartTime(e.target.value)}
                  className="pl-10"
                  disabled={isPending}
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="edit-end-time">End Time</Label>
              <div className="relative">
                <Clock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="edit-end-time"
                  type="time"
                  value={endTime}
                  onChange={(e) => setEndTime(e.target.value)}
                  className="pl-10"
                  disabled={isPending}
                />
              </div>
            </div>
          </div>

          <div className="flex items-center justify-between p-4 rounded-lg border">
            <div className="flex items-center gap-2">
              <Video className="h-4 w-4 text-muted-foreground" />
              <Label htmlFor="edit-is-online" className="cursor-pointer">
                Online Session
              </Label>
            </div>
            <Switch
              id="edit-is-online"
              checked={isOnline}
              onCheckedChange={setIsOnline}
              disabled={isPending}
            />
          </div>

          {isOnline ? (
            <div className="space-y-2">
              <Label htmlFor="edit-meeting-url">Meeting URL</Label>
              <Input
                id="edit-meeting-url"
                type="url"
                placeholder="https://zoom.us/j/..."
                value={meetingUrl}
                onChange={(e) => setMeetingUrl(e.target.value)}
                disabled={isPending}
              />
            </div>
          ) : (
            <div className="space-y-2">
              <Label htmlFor="edit-location">Location</Label>
              <div className="relative">
                <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="edit-location"
                  placeholder="e.g., Room 101"
                  value={location}
                  onChange={(e) => setLocation(e.target.value)}
                  className="pl-10"
                  disabled={isPending}
                />
              </div>
            </div>
          )}

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
                  <AlertDialogTitle>Delete Lecture</AlertDialogTitle>
                  <AlertDialogDescription>
                    Are you sure you want to delete this lecture? This action
                    cannot be undone.
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
