"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Upload, Loader2, FileText, Link as LinkIcon } from "lucide-react";

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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { submissionService } from "@/services/submissions";

const submitAssignmentSchema = z
  .object({
    content_text: z.string().max(50000).optional(),
    file_url: z
      .string()
      .url("Must be a valid URL")
      .max(500)
      .optional()
      .or(z.literal("")),
  })
  .refine((data) => data.content_text || data.file_url, {
    message: "Please provide either text content or a file URL",
  });

type FormData = z.infer<typeof submitAssignmentSchema>;

interface SubmitAssignmentDialogProps {
  assignmentId: string;
  assignmentTitle: string;
  trigger?: React.ReactNode;
  onSuccess?: () => void;
}

export function SubmitAssignmentDialog({
  assignmentId,
  assignmentTitle,
  trigger,
  onSuccess,
}: SubmitAssignmentDialogProps) {
  const [open, setOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<"text" | "file">("text");
  const queryClient = useQueryClient();

  const form = useForm<FormData>({
    resolver: zodResolver(submitAssignmentSchema),
    defaultValues: {
      content_text: "",
      file_url: "",
    },
  });

  const mutation = useMutation({
    mutationFn: (data: FormData) =>
      submissionService.create({
        assignment_id: assignmentId,
        content_text: data.content_text || undefined,
        file_url: data.file_url || undefined,
      }),
    onSuccess: () => {
      toast.success("Assignment submitted successfully!");
      queryClient.invalidateQueries({ queryKey: ["submissions"] });
      queryClient.invalidateQueries({ queryKey: ["my-submissions"] });
      queryClient.invalidateQueries({ queryKey: ["assignment", assignmentId] });
      setOpen(false);
      form.reset();
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to submit assignment");
    },
  });

  const onSubmit = (data: FormData) => {
    mutation.mutate(data);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button>
            <Upload className="mr-2 h-4 w-4" />
            Submit Assignment
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-125">
        <DialogHeader>
          <DialogTitle>Submit Assignment</DialogTitle>
          <DialogDescription>Submitting: {assignmentTitle}</DialogDescription>
        </DialogHeader>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <Tabs
            value={activeTab}
            onValueChange={(v) => setActiveTab(v as "text" | "file")}
          >
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="text" className="gap-2">
                <FileText className="h-4 w-4" />
                Text
              </TabsTrigger>
              <TabsTrigger value="file" className="gap-2">
                <LinkIcon className="h-4 w-4" />
                File URL
              </TabsTrigger>
            </TabsList>

            <TabsContent value="text" className="space-y-4 mt-4">
              <div className="space-y-2">
                <Label htmlFor="content_text">Your Answer</Label>
                <Textarea
                  id="content_text"
                  placeholder="Type your answer here..."
                  className="resize-none min-h-50 font-mono"
                  {...form.register("content_text")}
                />
                <p className="text-xs text-muted-foreground">
                  Write your solution or answer directly here
                </p>
              </div>
            </TabsContent>

            <TabsContent value="file" className="space-y-4 mt-4">
              <div className="space-y-2">
                <Label htmlFor="file_url">File URL</Label>
                <Input
                  id="file_url"
                  placeholder="https://drive.google.com/file/..."
                  {...form.register("file_url")}
                />
                <p className="text-xs text-muted-foreground">
                  Paste a link to your file (Google Drive, Dropbox, GitHub,
                  etc.)
                </p>
                {form.formState.errors.file_url && (
                  <p className="text-sm text-destructive">
                    {form.formState.errors.file_url.message}
                  </p>
                )}
              </div>
              <div className="rounded-lg border-2 border-dashed p-6 text-center">
                <Upload className="mx-auto h-8 w-8 text-muted-foreground mb-2" />
                <p className="text-sm text-muted-foreground">
                  Upload your file to a cloud service and paste the sharing link
                  above
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                  Supported: Google Drive, Dropbox, OneDrive, GitHub
                </p>
              </div>
            </TabsContent>
          </Tabs>

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
              Submit
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
