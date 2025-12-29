import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { execSync } from 'child_process';

let buildVersion;
try {
    // %h = short hash
    // %cd = committer date
    // --date=format:... = custom date format
    buildVersion = execSync('git log -1 --format="%h-%cd" --date=format:"%Y%m%d-%H"').toString().trim();
} catch (e) {
    buildVersion = 'dev-' + new Date().toISOString().replace(/[:.]/g, '-').slice(0, 13);
}

export default defineConfig({
    plugins: [tailwindcss(), sveltekit()],
    define: {
        __BUILD_VERSION__: JSON.stringify(buildVersion)
    } });
