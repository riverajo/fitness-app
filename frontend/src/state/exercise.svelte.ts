import { type Client } from '@urql/svelte';
import { graphql } from '$lib/gql';
import { browser } from '$app/environment';

export interface Exercise {
	id: string;
	name: string;
	description?: string | null;
	isCustom: boolean;
}

const GET_ALL_EXERCISES = graphql(`
	query GetAllExercises {
		uniqueExercises(limit: 1000) {
			id
			name
			description
			isCustom
		}
	}
`);

export class ExerciseStore {
	#exercises = $state<Exercise[]>([]);
	#loading = $state(false);
	#initialized = $state(false);

	constructor() {
		if (browser) {
			this.loadFromStorage();
		}
	}

	get all() {
		return this.#exercises;
	}

	get loading() {
		return this.#loading;
	}

	get initialized() {
		return this.#initialized;
	}

	loadFromStorage() {
		const stored = localStorage.getItem('exercises_cache');
		if (stored) {
			try {
				const parased = JSON.parse(stored);
				this.#exercises = parased;
				this.#initialized = true;
			} catch (e) {
				console.error('Failed to parse exercises from local storage', e);
			}
		}
	}

	async sync(client: Client) {
		this.#loading = true;
		try {
			const result = await client.query(GET_ALL_EXERCISES, {}).toPromise();

			if (result.data?.uniqueExercises) {
				this.#exercises = result.data.uniqueExercises;
				if (browser) {
					localStorage.setItem('exercises_cache', JSON.stringify(this.#exercises));
				}
			}
		} catch (e) {
			console.error('Failed to sync exercises', e);
		} finally {
			this.#loading = false;
			this.#initialized = true;
		}
	}

	search(query: string): Exercise[] {
		if (!query) return this.#exercises;
		const lowerQuery = query.toLowerCase();
		return this.#exercises.filter((exercise) => exercise.name.toLowerCase().includes(lowerQuery));
	}
}

export const exerciseStore = new ExerciseStore();
