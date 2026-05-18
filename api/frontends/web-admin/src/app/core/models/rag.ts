export interface RagQueryRequest {
  question: string;
  top_k?: number;
}

export interface RagQueryResponse {
  answer: string;
  document_ids: string[];
  snippets: string[];
}

export interface RagDocumentUploadResponse {
  document_id: string;
  status: string;
  chunk_count: number;
}

export interface RagDocumentDeleteResponse {
  document_id: string;
  deleted: boolean;
}
