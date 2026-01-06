<script lang="ts">
	import { Input } from 'flowbite-svelte';

	export let value: number = 0;
	export let size: 'sm' | 'md' | 'lg' = 'sm';

	let lbs = 0;
	let kgs = value;

	// Constant for conversion
	const LBS_TO_KG = 0.45359237;

	// Track the last value we emitted to detecting external changes
	let lastEmittedValue = value;

	// Watch for external changes to value (e.g. loading from DB or reset form)
	$: if (value !== lastEmittedValue) {
		// If the value changes from outside, we reset our local inputs
		// Default to all KGs since we don't store the split
		kgs = value;
		lbs = 0;
		lastEmittedValue = value;
	}

	function updateValue(_triggerLbs?: number, _triggerKgs?: number) {
		let total = lbs * LBS_TO_KG + kgs;
		// Round to 2 decimal places
		total = Math.round(total * 100) / 100;

		lastEmittedValue = total;
		value = total;
	}

	// Trigger update whenever inputs change (replaces on:input to avoid binding race conditions)
	$: updateValue(lbs, kgs);
</script>

<div class="flex items-center gap-2">
	<div class="flex flex-col">
		<div class="flex items-center gap-1">
			<Input
				type="number"
				{size}
				class="w-20"
				bind:value={lbs}
				placeholder="0"
				data-testid="weight-input-lbs"
			/>
			<span class="text-xs text-gray-500">lbs</span>
		</div>
	</div>
	<span class="text-gray-400">+</span>
	<div class="flex flex-col">
		<div class="flex items-center gap-1">
			<Input
				type="number"
				{size}
				class="w-20"
				bind:value={kgs}
				placeholder="0"
				data-testid="weight-input-kgs"
			/>
			<span class="text-xs text-gray-500">kgs</span>
		</div>
	</div>
</div>
