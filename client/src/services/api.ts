import apiHandler from '@/handlers/api-handlers';
import { Post, Comment, PostsResponse, CommentsResponse, AuthResponse, SecureAuthResponse, LoginRequest, RegisterRequest } from '@/types';
import { toast } from 'sonner';

export class AuthAPI {
  static async login(credentials: LoginRequest): Promise<AuthResponse> {
    try {
      const response = await apiHandler.post<AuthResponse>('/auth/login', credentials);
      if (response.success) {
        // Store tokens and user data
        localStorage.setItem('access_token', response.data!.access_token);
        localStorage.setItem('refresh_token', response.data!.refresh_token);
        localStorage.setItem('user', JSON.stringify(response.data!.user));
        return response.data!;
      } else {
        toast.error(response.error_message || 'Login failed');
        throw new Error(response.error_message || 'Login failed');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Login failed';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async register(userData: RegisterRequest): Promise<AuthResponse> {
    try {
      const response = await apiHandler.post<AuthResponse>('/auth/register', userData);
      if (response.success) {
        // Store tokens and user data
        localStorage.setItem('access_token', response.data!.access_token);
        localStorage.setItem('refresh_token', response.data!.refresh_token);
        localStorage.setItem('user', JSON.stringify(response.data!.user));
        return response.data!;
      } else {
        toast.error(response.error_message || 'Registration failed');
        throw new Error(response.error_message || 'Registration failed');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Registration failed';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async logout(): Promise<void> {
    try {
      await apiHandler.post('/auth/logout');
    } catch (error) {
      // Continue with logout even if API call fails
      console.error('Logout API call failed:', error);
    } finally {
      // Always clear local storage
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user');
    }
  }

  static async refreshToken(): Promise<SecureAuthResponse> {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    try {
      const response = await apiHandler.post<SecureAuthResponse>('/auth/refresh', {
        refresh_token: refreshToken
      });
      if (response.success) {
        // Update stored access token and user data (refresh token stays the same)
        localStorage.setItem('access_token', response.data!.access_token);
        localStorage.setItem('user', JSON.stringify(response.data!.user));
        return response.data!;
      } else {
        throw new Error(response.error_message || 'Token refresh failed');
      }
    } catch (error: any) {
      // Clear tokens on refresh failure
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user');
      throw error;
    }
  }

  static getCurrentUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  }

  static getAccessToken() {
    return localStorage.getItem('access_token');
  }

  static isAuthenticated() {
    return !!localStorage.getItem('access_token');
  }
}

export class PostsAPI {
  static async getPosts(page: number = 1, limit: number = 10): Promise<PostsResponse> {
    try {
      const offset = (page - 1) * limit;
      const response = await apiHandler.get<any>(`/posts?limit=${limit}&offset=${offset}`);
      if (response.success) {
        // Transform backend response to frontend format
        const backendData = response.data!;
        const transformedResponse: PostsResponse = {
          posts: backendData.posts || [],
          meta: {
            page: page,
            limit: limit,
            total: backendData.count || 0,
            totalPages: Math.ceil((backendData.count || 0) / limit),
            hasNext: (backendData.count || 0) > offset + limit,
            hasPrev: page > 1
          }
        };
        return transformedResponse;
      } else {
        toast.error(response.error_message || 'Failed to fetch posts');
        throw new Error(response.error_message || 'Failed to fetch posts');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to fetch posts';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async getPost(postId: string): Promise<Post> {
    try {
      const response = await apiHandler.get<Post>(`/posts/post/${postId}`);
      if (response.success) {
        return response.data!;
      } else {
        toast.error(response.error_message || 'Failed to fetch post');
        throw new Error(response.error_message || 'Failed to fetch post');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to fetch post';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async createPost(postData: { title: string; content: string }): Promise<Post> {
    try {
      const response = await apiHandler.post<Post>('/posts', postData);
      if (response.success) {
        toast.success('Post created successfully');
        return response.data!;
      } else {
        toast.error(response.error_message || 'Failed to create post');
        throw new Error(response.error_message || 'Failed to create post');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to create post';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async updatePost(postId: string, postData: { title?: string; content?: string }): Promise<Post> {
    try {
      const response = await apiHandler.put<Post>(`/posts/post/${postId}`, postData);
      if (response.success) {
        toast.success('Post updated successfully');
        return response.data!;
      } else {
        toast.error(response.error_message || 'Failed to update post');
        throw new Error(response.error_message || 'Failed to update post');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to update post';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async deletePost(postId: string): Promise<void> {
    try {
      const response = await apiHandler.delete(`/posts/post/${postId}`);
      if (response.success) {
        toast.success('Post deleted successfully');
      } else {
        toast.error(response.error_message || 'Failed to delete post');
        throw new Error(response.error_message || 'Failed to delete post');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to delete post';
      toast.error(errorMessage);
      throw error;
    }
  }
}

export class CommentsAPI {
  static async getComments(postId: string, page: number = 1, limit: number = 10): Promise<CommentsResponse> {
    try {
      const offset = (page - 1) * limit;
      const response = await apiHandler.get<any>(`/posts/post-comments/${postId}?limit=${limit}&offset=${offset}`);
      if (response.success) {
        // Transform backend response to frontend format
        const backendData = response.data!;
        const transformedResponse: CommentsResponse = {
          comments: backendData.comments || [],
          meta: {
            page: page,
            limit: limit,
            total: backendData.count || 0,
            totalPages: Math.ceil((backendData.count || 0) / limit),
            hasNext: (backendData.count || 0) > offset + limit,
            hasPrev: page > 1
          }
        };
        return transformedResponse;
      } else {
        toast.error(response.error_message || 'Failed to fetch comments');
        throw new Error(response.error_message || 'Failed to fetch comments');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to fetch comments';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async getReplies(commentId: string, page: number = 1, limit: number = 5): Promise<CommentsResponse> {
    try {
      const offset = (page - 1) * limit;
      const response = await apiHandler.get<any>(`/comments/${commentId}/replies?limit=${limit}&offset=${offset}`);
      if (response.success) {
        // Transform backend response to frontend format
        const backendData = response.data!;
        const transformedResponse: CommentsResponse = {
          comments: backendData.replies || [],
          meta: {
            page: page,
            limit: limit,
            total: backendData.count || 0,
            totalPages: Math.ceil((backendData.count || 0) / limit),
            hasNext: (backendData.count || 0) > offset + limit,
            hasPrev: page > 1
          }
        };
        return transformedResponse;
      } else {
        toast.error(response.error_message || 'Failed to fetch replies');
        throw new Error(response.error_message || 'Failed to fetch replies');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to fetch replies';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async createComment(postId: string, content: string, parentId?: string): Promise<Comment> {
    try {
      const commentData = {
        content,
        post_id: postId,
        ...(parentId && { parent_id: parentId })
      };
      const response = await apiHandler.post<Comment>(`/posts/post-comments/${postId}`, commentData);
      if (response.success) {
        toast.success('Comment added successfully');
        return response.data!;
      } else {
        toast.error(response.error_message || 'Failed to add comment');
        throw new Error(response.error_message || 'Failed to add comment');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to add comment';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async updateComment(commentId: string, content: string): Promise<Comment> {
    try {
      const response = await apiHandler.put<Comment>(`/comments/${commentId}`, { content });
      if (response.success) {
        toast.success('Comment updated successfully');
        return response.data!;
      } else {
        toast.error(response.error_message || 'Failed to update comment');
        throw new Error(response.error_message || 'Failed to update comment');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to update comment';
      toast.error(errorMessage);
      throw error;
    }
  }

  static async deleteComment(commentId: string): Promise<void> {
    try {
      const response = await apiHandler.delete(`/comments/${commentId}`);
      if (response.success) {
        toast.success('Comment deleted successfully');
      } else {
        toast.error(response.error_message || 'Failed to delete comment');
        throw new Error(response.error_message || 'Failed to delete comment');
      }
    } catch (error: any) {
      const errorMessage = error.message || 'Failed to delete comment';
      toast.error(errorMessage);
      throw error;
    }
  }
} 