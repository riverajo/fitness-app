<script lang="ts">
    import { gql, getContextClient, mutationStore } from '@urql/svelte';
    import { goto } from '$app/navigation';

    const client = getContextClient();

    let name = '';
    let description = '';
    let error = '';

    const createExerciseMutation = gql`
        mutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {
            createUniqueExercise(input: $input) {
                id
                name
                description
                isCustom
            }
        }
    `;

    async function handleSubmit() {
        error = '';
        const result = await client.mutation(createExerciseMutation, {
            input: {
                name,
                description
            }
        }).toPromise();

        if (result.error) {
            error = result.error.message;
        } else {
            goto('/exercises');
        }
    }
</script>

<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
    <h1 class="text-3xl font-bold tracking-tight text-gray-900">Create New Exercise</h1>

    <form on:submit|preventDefault={handleSubmit} class="mt-6 space-y-6 max-w-xl">
        {#if error}
            <div class="rounded-md bg-red-50 p-4">
                <div class="flex">
                    <div class="ml-3">
                        <h3 class="text-sm font-medium text-red-800">Error creating exercise</h3>
                        <div class="mt-2 text-sm text-red-700">
                            <p>{error}</p>
                        </div>
                    </div>
                </div>
            </div>
        {/if}

        <div>
            <label for="name" class="block text-sm font-medium leading-6 text-gray-900">Name</label>
            <div class="mt-2">
                <input
                    type="text"
                    name="name"
                    id="name"
                    required
                    bind:value={name}
                    class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                    placeholder="e.g. Bench Press"
                />
            </div>
        </div>

        <div>
            <label for="description" class="block text-sm font-medium leading-6 text-gray-900">Description</label>
            <div class="mt-2">
                <textarea
                    id="description"
                    name="description"
                    rows="3"
                    bind:value={description}
                    class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                    placeholder="Optional description..."
                ></textarea>
            </div>
        </div>

        <div class="flex items-center justify-end gap-x-6">
            <a href="/dashboard" class="text-sm font-semibold leading-6 text-gray-900">Cancel</a>
            <button
                type="submit"
                class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
                Create Exercise
            </button>
        </div>
    </form>
</div>
