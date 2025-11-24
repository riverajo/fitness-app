---
description: Context and patterns for AI agents working on this project.
---
# AI Context: Flowbite-Svelte

## Overview
We use **Flowbite-Svelte** with **Tailwind CSS v4**.

## Documentation
- **Official Docs**: [https://flowbite-svelte.com/](https://flowbite-svelte.com/)
- **LLM Docs**: [https://flowbite-svelte.com/llms.txt](https://flowbite-svelte.com/llms.txt) (Use this for RAG if available)

## Common Patterns

### Buttons
```svelte
<script>
  import { Button } from 'flowbite-svelte';
</script>

<Button color="blue">Click me</Button>
<Button outline>Outline</Button>
```

### Bottom Navigation (Mobile)
```svelte
<script>
  import { BottomNav, BottomNavItem } from 'flowbite-svelte';
</script>

<BottomNav position="absolute" classInner="grid-cols-3">
  <BottomNavItem btnName="Home" href="/">
    <!-- Icon here -->
  </BottomNavItem>
  <BottomNavItem btnName="Workouts" href="/workouts">
    <!-- Icon here -->
  </BottomNavItem>
</BottomNav>
```

### Drawers (Mobile Menus)
```svelte
<script>
  import { Drawer, Button, CloseButton } from 'flowbite-svelte';
  import { sineIn } from 'svelte/easing';
  let hidden = true;
  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn
  };
</script>

<Drawer transitionType="fly" {transitionParams} bind:hidden>
  <CloseButton on:click={() => (hidden = true)} />
  <!-- Content -->
</Drawer>
```

## Tailwind v4 Note
We use `@source "../node_modules/flowbite-svelte";` in our CSS. You do not need to configure `tailwind.config.js` content paths manually.
