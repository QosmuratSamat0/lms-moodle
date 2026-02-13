"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  CalendarIcon,
  Plus,
  Loader2,
  Clock,
  MapPin,
  Video,
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
import type { EventType, RecurrenceType } from "@/types/schedule";

interface CreateLectureDialogProps {
  courseId: string;
  trigger?: React.ReactNode;
}

export function CreateLectureDialog({
  courseId,
  trigger,
}: CreateLectureDialogProps) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [eventType, setEventType] = useState<EventType>("class");
  const [startDate, setStartDate] = useState<Date | undefined>();
  const [startTime, setStartTime] = useState("09:00");
  const [endTime, setEndTime] = useState("10:30");
  const [location, setLocation] = useState("");
  const [isOnline, setIsOnline] = useState(false);
  const [meetingUrl, setMeetingUrl] = useState("");
  const [recurrence, setRecurrence] = useState<RecurrenceType>("none");
  const [errors, setErrors] = useState<{ title?: string; date?: string }>({});
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: () => {
      if (!startDate) throw new Error("Start date is required");

      const startDateTime = new Date(startDate);
      const [startHour, startMinute] = startTime.split(":").map(Number);
      startDateTime.setHours(startHour, startMinute, 0, 0);

      const endDateTime = new Date(startDate);
      const [endHour, endMinute] = endTime.split(":").map(Number);
      endDateTime.setHours(endHour, endMinute, 0, 0);

      return scheduleService.create({
        course_id: courseId,
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
      toast.success("Lecture created successfully!");
      queryClient.invalidateQueries({
        queryKey: ["course-schedule", courseId],
      });
      queryClient.invalidateQueries({ queryKey: ["schedule"] });
      setOpen(false);
      resetForm();
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to create lecture");
    },
  });

  const resetForm = () => {
    setTitle("");
    setDescription("");
    setEventType("class");
    setStartDate(undefined);
    setStartTime("09:00");
    setEndTime("10:30");
    setLocation("");
    setIsOnline(false);
    setMeetingUrl("");
    setRecurrence("none");
    setErrors({});
  };

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
            Add Lecture
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-137.5 max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Add New Lecture</DialogTitle>
          <DialogDescription>
            Create a new lecture, lab session, or exam for this course.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title *</Label>
            <Input
              id="title"
              placeholder="e.g., Introduction to Programming"
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                if (errors.title)
                  setErrors((prev) => ({ ...prev, title: undefined }));
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
              placeholder="Lecture description or topics to cover..."
              className="resize-none min-h-20"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Event Type</Label>
              <Select
                value={eventType}
                onValueChange={(v) => setEventType(v as EventType)}
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
              <Label htmlFor="start-time">Start Time</Label>
              <div className="relative">
                <Clock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="start-time"
                  type="time"
                  value={startTime}
                  onChange={(e) => setStartTime(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="end-time">End Time</Label>
              <div className="relative">
                <Clock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="end-time"
                  type="time"
                  value={endTime}
                  onChange={(e) => setEndTime(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
          </div>

          <div className="flex items-center justify-between p-4 rounded-lg border">
            <div className="flex items-center gap-2">
              <Video className="h-4 w-4 text-muted-foreground" />
              <Label htmlFor="is-online" className="cursor-pointer">
                Online Session
              </Label>
            </div>
            <Switch
              id="is-online"
              checked={isOnline}
              onCheckedChange={setIsOnline}
            />
          </div>

          {isOnline ? (
            <div className="space-y-2">
              <Label htmlFor="meeting-url">Meeting URL</Label>
              <Input
                id="meeting-url"
                type="url"
                placeholder="https://zoom.us/j/..."
                value={meetingUrl}
                onChange={(e) => setMeetingUrl(e.target.value)}
              />
            </div>
          ) : (
            <div className="space-y-2">
              <Label htmlFor="location">Location</Label>
              <div className="relative">
                <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  id="location"
                  placeholder="e.g., Room 101"
                  value={location}
                  onChange={(e) => setLocation(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
          )}

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
              Create Lecture
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
