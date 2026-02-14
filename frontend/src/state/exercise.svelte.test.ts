import { describe, it, expect, beforeEach, vi } from 'vitest';
import type { Client } from '@urql/svelte';
import { ExerciseStore } from './exercise.svelte';

// Mock $app/environment
vi.mock('$app/environment', () => ({
	browser: true
}));

// Mock localStorage
const localStorageMock = (function () {
	let store: Record<string, string> = {};
	return {
		getItem: function (key: string) {
			return store[key] || null;
		},
		setItem: function (key: string, value: string) {
			store[key] = value.toString();
		},
		clear: function () {
			store = {};
		},
		removeItem: function (key: string) {
			delete store[key];
		}
	};
})();

Object.defineProperty(window, 'localStorage', {
	value: localStorageMock
});

describe('ExerciseStore', () => {
	beforeEach(() => {
		localStorage.clear();
		vi.clearAllMocks();
	});

	it('should initialize with empty exercises if storage is empty', () => {
		const store = new ExerciseStore();
		expect(store.all).toEqual([]);
		expect(store.loading).toBe(false);
	});

	it('should load exercises from local storage on init', () => {
		const mockExercises = [
			{ id: '1', name: 'Bench Press', isCustom: false },
			{ id: '2', name: 'Squat', isCustom: false }
		];
		localStorage.setItem('exercises_cache', JSON.stringify(mockExercises));

		const store = new ExerciseStore();
		expect(store.all).toEqual(mockExercises);
		expect(store.initialized).toBe(true);
	});

	it('should sync exercises from client', async () => {
		const mockExercises = [
			{ id: '1', name: 'Bench Press', isCustom: false },
			{ id: '2', name: 'Squat', isCustom: false }
		];

		const mockClient = {
			query: vi.fn().mockReturnValue({
				toPromise: vi.fn().mockResolvedValue({
					data: { uniqueExercises: mockExercises }
				})
			})
		};

		const store = new ExerciseStore();
		await store.sync(mockClient as unknown as Client);

		expect(store.all).toEqual(mockExercises);
		expect(store.initialized).toBe(true);
		expect(localStorage.getItem('exercises_cache')).toBe(JSON.stringify(mockExercises));
	});

	it('should filter exercises by search query', () => {
		const mockExercises = [
			{ id: '1', name: 'Bench Press', isCustom: false },
			{ id: '2', name: 'Squat', isCustom: false },
			{ id: '3', name: 'Overhead Press', isCustom: false }
		];
		localStorage.setItem('exercises_cache', JSON.stringify(mockExercises));

		const store = new ExerciseStore();

		const results = store.search('Press');
		expect(results).toHaveLength(2);
		expect(results.map((e) => e.name)).toContain('Bench Press');
		expect(results.map((e) => e.name)).toContain('Overhead Press');
		expect(results.map((e) => e.name)).not.toContain('Squat');
	});

	it('should handle empty search query', () => {
		const mockExercises = [{ id: '1', name: 'Bench Press', isCustom: false }];
		localStorage.setItem('exercises_cache', JSON.stringify(mockExercises));

		const store = new ExerciseStore();
		const results = store.search('');
		expect(results).toEqual(mockExercises);
	});
});
