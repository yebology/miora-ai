// Matches backend/app/output/response.go ApiResponse envelope
export type ApiResponse<T> = {
  status: "success" | "error";
  message: string;
  data?: T;
};
