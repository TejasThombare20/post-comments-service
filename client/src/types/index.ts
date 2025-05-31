export interface User {
  id: string;
  username: string;
  email: string;
  display_name: string;
  avatar_url: string | null;
  created_at: string;
  updated_at: string;
}

export interface Post {
  id: string;
  title: string;
  content: string;
  created_by: string;
  author: User;
  created_at: string;
  updated_at: string;
  commentsCount?: number;
}

export interface Comment {
  id: string;
  content: string;
  post_id: string;
  parent_id?: string;
  path: string[];
  thread_id: string;
  created_by?: string;
  author?: User;
  children?: Comment[];
  created_at: string;
  updated_at: string;
  replies_count: number;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface PostsResponse {
  posts: Post[];
  meta: PaginationMeta;
}

export interface CommentsResponse {
  comments: Comment[];
  meta: PaginationMeta;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_at: string;
}

export interface SecureAuthResponse {
  user: User;
  access_token: string;
  expires_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  display_name: string;
} 