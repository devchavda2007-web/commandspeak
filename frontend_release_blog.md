# The Face of CommandSpeak: Building an Offline-First Web UI

During the **CYHI Hackathon (Track 4)**, the core challenge was creating a tool that removed friction from developer workflows. We solved the backend logic with our fast, local keyword-scoring CLI engine. But we knew we could push it further.

What if you could build your scripts visually, check your commands in a simulated terminal, and export them without opening a terminal at all?

Enter the **CommandSpeak Web UI** (`index.html`).

### Why a Web UI?
Command-line interfaces are fast, but they lack discoverability. The Web UI acts as an interactive playground where developers can speak or type natural language and instantly see the translated shell command. 

### Key Features Implemented:
1. **100% Offline PWA**: We converted the UI into a Progressive Web App (PWA). By adding a `manifest.json` and a Service Worker (`sw.js`), developers can install CommandSpeak on their desktops or mobile devices. Because our NLP engine relies on keyword-scoring rather than LLM APIs, it runs completely offline in the browser.
2. **Export as Shell Script (.sh)**: Tracking command history is nice, but executing it is better. We added a feature to let developers export their session history directly into a `.sh` bash script, effectively turning the UI into a visual script-builder.
3. **Voice to Code**: Using the Web Speech API, developers can literally talk to their terminal UI. Click the mic, say "deploy my project to vercel", and the command is generated.
4. **Theme & Accessibility**: We baked in a light/dark mode toggle and ensured our custom toggles use proper ARIA labels so the tool is accessible to all developers.

### The Result
The combination of a powerful Go-based CLI and an offline-first PWA frontend creates a complete ecosystem. CommandSpeak isn't just a tool; it's a new way to interact with your workspace.

*(Log generated during the final hours of the 24-hour CYHI Hackathon)*
