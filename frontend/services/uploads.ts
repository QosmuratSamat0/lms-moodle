import api from "@/lib/api-client";

export interface UploadedFile {
  id: string;
  uploaded_by: string;
  type: string;
  reference_id?: string;
  public_id: string;
  url: string;
  secure_url: string;
  original_name: string;
  format: string;
  resource_type: string;
  size: number;
  width?: number;
  height?: number;
  uploader_first_name?: string;
  uploader_last_name?: string;
  uploader_email?: string;
  created_at: string;
}

// Alias for convenience
export interface UploadResponse {
  file_url: string;
  file_name: string;
}

export type AttachmentType =
  | "submission"
  | "assignment"
  | "course"
  | "profile"
  | "chat";

export const uploadService = {
  /**
   * Upload a file to Cloudinary via the backend
   * @param file - The file to upload
   * @param type - The attachment type (submission, assignment, course, profile, chat)
   * @param referenceId - Optional reference ID (submission_id, assignment_id, etc.)
   */
  upload: async (
    file: File,
    type: AttachmentType,
    referenceId?: string,
  ): Promise<UploadedFile> => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("type", type);
    if (referenceId) {
      formData.append("reference_id", referenceId);
    }

    const result = await api.upload<UploadedFile>("/uploads", formData);

    // Add convenience properties to match what submission expects
    return {
      ...result,
      // These are for convenience - map secure_url to file_url
      file_url: result.secure_url || result.url,
      file_name: result.original_name,
    } as UploadedFile & { file_url: string; file_name: string };
  },

  /**
   * Get uploads by reference (e.g., all files for a submission)
   */
  getByReference: async (
    type: AttachmentType,
    referenceId: string,
  ): Promise<UploadedFile[]> => {
    return api.get<UploadedFile[]>(`/uploads/reference/${type}/${referenceId}`);
  },

  /**
   * Get a single upload by ID
   */
  get: async (id: string): Promise<UploadedFile> => {
    return api.get<UploadedFile>(`/uploads/${id}`);
  },

  /**
   * Delete an upload
   */
  delete: async (id: string): Promise<void> => {
    return api.delete(`/uploads/${id}`);
  },

  /**
   * Get all uploads for the current user
   */
  getMyUploads: async (): Promise<UploadedFile[]> => {
    return api.get<UploadedFile[]>("/uploads/my");
  },

  /**
   * Get allowed file types
   */
  getAllowedTypes: async (): Promise<{
    extensions: string[];
    max_file_size: number;
  }> => {
    return api.get<{ extensions: string[]; max_file_size: number }>(
      "/uploads/allowed-types",
    );
  },
};

export default uploadService;
