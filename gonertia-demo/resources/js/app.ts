import { createInertiaApp } from '@inertiajs/svelte'

createInertiaApp({
  resolve: name => {
    // @ts-expect-error
    const pages = import.meta.glob("./Pages/**/*.svelte", { eager: true });
    return pages[`./Pages/${name}.svelte`];
  },
  setup({ el, App, props }) {
    new App({ target: el, props })
  },
})
