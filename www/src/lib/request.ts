import type { Response, PaginatedResponse } from './types';
import { notifyError } from './stores/notificationStore';

/**
 * Custom API request options
 */
export interface RequestOptions extends RequestInit {
	handleRedirects?: boolean;
	notifyOnError?: boolean;
}

/**
 * Generic API request function with standardized error handling
 * @param url The URL to fetch
 * @param options Request options
 * @returns Promise with the API response
 */
export async function apiRequest<T>(
	url: string,
	options: RequestOptions = {}
): Promise<Response<T>> {
	// Default options
	const defaultOptions: RequestOptions = {
		credentials: 'include', // Include cookies by default
		redirect: options.handleRedirects === false ? 'manual' : 'follow',
		notifyOnError: true,
		...options
	};

	try {
		const response = await fetch(url, defaultOptions);

		// Handle unauthorized access (401)
		if (response.status === 401) {
			const data = (await response.json()) as Response<null>;

			// If the backend provides a redirect URL, use it
			if (data.error?.redirect) {
				window.location.href = data.error.redirect;
				return { meta: { success: false }, data: null, error: data.error } as Response<T>;
			} else {
				// Default unauthorized handling
				const error = { message: 'You need to log in to access this resource' };
				if (defaultOptions.notifyOnError) {
					notifyError(error.message);
				}
				return { meta: { success: false }, data: null, error } as Response<T>;
			}
		}

		// Handle other errors
		if (!response.ok) {
			const errorData = await response.json().catch(() => null);
			const error = errorData?.error || {
				message: `Request failed with status: ${response.status}`
			};

			if (defaultOptions.notifyOnError) {
				notifyError(error.message);
			}

			return { meta: { success: false }, data: null, error } as Response<T>;
		}

		// Handle successful responses
		const data = await response.json();
		return data as Response<T>;
	} catch (error) {
		const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';

		if (defaultOptions.notifyOnError) {
			notifyError(errorMessage);
		}

		return {
			meta: { success: false },
			data: null,
			error: { message: errorMessage }
		} as Response<T>;
	}
}

/**
 * GET request helper
 */
export function get<T>(url: string, options: RequestOptions = {}): Promise<Response<T>> {
	return apiRequest<T>(url, { ...options, method: 'GET' });
}

/**
 * POST request helper
 */
export function post<T>(
	url: string,
	data?: any,
	options: RequestOptions = {}
): Promise<Response<T>> {
	const headers = {
		'Content-Type': 'application/json',
		...options.headers
	};

	return apiRequest<T>(url, {
		...options,
		method: 'POST',
		headers,
		body: data ? JSON.stringify(data) : undefined
	});
}

/**
 * PUT request helper
 */
export function put<T>(
	url: string,
	data?: any,
	options: RequestOptions = {}
): Promise<Response<T>> {
	const headers = {
		'Content-Type': 'application/json',
		...options.headers
	};

	return apiRequest<T>(url, {
		...options,
		method: 'PUT',
		headers,
		body: data ? JSON.stringify(data) : undefined
	});
}

/**
 * DELETE request helper
 */
export function del<T>(url: string, options: RequestOptions = {}): Promise<Response<T>> {
	return apiRequest<T>(url, { ...options, method: 'DELETE' });
}

/**
 * Helper function for paginated API requests
 * @param url The URL to fetch
 * @param options Request options
 * @returns Promise with the paginated API response
 */
export async function getPaginated<T>(
	url: string,
	options: RequestOptions = {}
): Promise<PaginatedResponse<T>> {
	try {
		const response = await fetch(url, {
			credentials: 'include',
			...options,
			method: 'GET'
		});

		if (!response.ok) {
			const errorData = await response.json().catch(() => null);
			const error = errorData?.error || {
				message: `Request failed with status: ${response.status}`
			};

			if (options.notifyOnError !== false) {
				notifyError(error.message);
			}

			throw new Error(error.message);
		}

		return (await response.json()) as PaginatedResponse<T>;
	} catch (error) {
		const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';

		if (options.notifyOnError !== false) {
			notifyError(errorMessage);
		}

		throw error;
	}
}
