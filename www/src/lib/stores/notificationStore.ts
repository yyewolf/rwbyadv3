import { writable } from 'svelte/store';

export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface Notification {
	id: string;
	type: NotificationType;
	message: string;
	timeout?: number; // Time in ms until notification auto-dismisses
}

function createNotificationStore() {
	const { subscribe, update } = writable<Notification[]>([]);

	const remove = (id: string) => {
		update((notifications) => {
			return notifications.filter((notification) => notification.id !== id);
		});
	};

	return {
		subscribe,
		add: (message: string, type: NotificationType = 'info', timeout: number = 5000) => {
			const id = Math.random().toString(36).substring(2, 9);

			update((notifications) => {
				return [...notifications, { id, type, message, timeout }];
			});

			if (timeout) {
				setTimeout(() => {
					remove(id);
				}, timeout);
			}

			return id;
		},
		remove: remove,
		clearAll: () => {
			update(() => []);
		}
	};
}

export const notificationStore = createNotificationStore();

// Helper functions for common notification types
export function notifySuccess(message: string, timeout = 5000) {
	return notificationStore.add(message, 'success', timeout);
}

export function notifyError(message: string, timeout = 8000) {
	return notificationStore.add(message, 'error', timeout);
}

export function notifyWarning(message: string, timeout = 5000) {
	return notificationStore.add(message, 'warning', timeout);
}

export function notifyInfo(message: string, timeout = 4000) {
	return notificationStore.add(message, 'info', timeout);
}
