<script lang="ts">
    import { mutationStore, gql, getContextClient } from '@urql/svelte';
    import { goto } from '$app/navigation';

    let email = $state('');
    let password = $state('');
    let error = $state('');
    let loading = $state(false);

    const registerMutation = gql`
        mutation Register($input: RegisterInput!) {
            register(input: $input) {
                success
                message
                user {
                    id
                    email
                }
            }
        }
    `;

    const client = getContextClient();

    async function handleSubmit(e: Event) {
        console.log('Submitting form...');
        e.preventDefault();
        loading = true;
        error = '';

        try {
            const result = await client.mutation(registerMutation, { input: { email, password } }).toPromise();

            if (result.error) {
                error = result.error.message;
            } else if (result.data?.register?.success) {
                goto('/');
            } else {
                error = result.data?.register?.message || 'Registration failed';
            }
        } catch (e) {
            error = 'An unexpected error occurred';
            console.error(e);
        } finally {
            loading = false;
        }
    }
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8">
    <div class="w-full max-w-md space-y-8">
        <div>
            <h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">Create your account</h2>
        </div>
        <form class="mt-8 space-y-6" onsubmit={handleSubmit}>
            <div class="-space-y-px rounded-md shadow-sm">
                <div>
                    <label for="email-address" class="sr-only">Email address</label>
                    <input id="email-address" name="email" type="email" autocomplete="email" required class="relative block w-full rounded-t-md border-0 py-1.5 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:z-10 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" placeholder="Email address" bind:value={email}>
                </div>
                <div>
                    <label for="password" class="sr-only">Password</label>
                    <input id="password" name="password" type="password" autocomplete="current-password" required class="relative block w-full rounded-b-md border-0 py-1.5 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:z-10 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" placeholder="Password" bind:value={password}>
                </div>
            </div>

            {#if error}
                <div class="rounded-md bg-red-50 p-4">
                    <div class="flex">
                        <div class="ml-3">
                            <h3 class="text-sm font-medium text-red-800">{error}</h3>
                        </div>
                    </div>
                </div>
            {/if}

            <div>
                <button type="submit" disabled={loading} class="group relative flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:opacity-50">
                    {#if loading}
                        Signing up...
                    {:else}
                        Sign up
                    {/if}
                </button>
            </div>
        </form>
    </div>
</div>
