import { createInertiaApp } from '@inertiajs/svelte'
import Layout from './Pages/Layout.svelte'

createInertiaApp({
  resolve: name => {
    // @ts-expect-error
    const pages = import.meta.glob("./Pages/**/*.svelte", { eager: true });
    let page = pages[`./Pages/${name}.svelte`];
    return { default: page.default, layout: page.layout || Layout };
  },
  setup({ el, App, props }) {
    new App({ target: el, props })
  },
})
