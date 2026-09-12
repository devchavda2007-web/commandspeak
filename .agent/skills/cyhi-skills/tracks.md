# CYHI - the four tracks

## Track 1 - Game: The 24-Hour Infinite Loop

Build an environmental game that pushes the limits of restricted key presses within a
relentless, resetting time cycle. Intuitive but deeply meaningful and enjoyable. Design a
tight survival or puzzle loop where the thrill comes entirely from the player's quick
thinking and perfectly timed execution.

**Hard constraint - allowed inputs only:**
- Swipe - a directional drag of the cursor
- WASD - standard movement keys
- Shift and Space - modifier and action keys

Mechanics, time-based puzzles, and survival objectives must be built entirely around those
three interactions. Nothing else. No mouse clicks, no other keys.

**Judges reward:** a tight, polished loop. One mechanic done well beats three half-finished
ones. The constraint is the design challenge - it should feel intentional, not limiting.

---

## Track 2 - Campus: Solve a Campus Problem

Every campus runs on workarounds: the WhatsApp group doing five jobs badly, the notice board
nobody checks, the question you ask five people before someone finally knows the answer.
Pick one and build software for it. Not a startup pitch - one real problem, solved well
enough that someone would actually use it.

**Objectives:**
- **Name the problem** - one specific problem, and say clearly who has it.
- **Ship it live** - deployed, openable by a real person tomorrow. Not a mockup.
- **Solve one thing** - a narrow tool that works beats a broad platform that half works.
- **Make it stick** - still works next month, not just during the demo.

**Judges reward:** specificity. The person sitting next to you at the demo should
immediately say "I would use this."

---

## Track 3 - AI/ML: A Layer Between You and the Noise

An AI-powered wellbeing layer sitting between an online community and the person reading it.
Catch harmful content before it reaches the reader, let people decide for themselves what
enters their feed, and show them what an hour of scrolling has done to their mood. Features
hand control back to the reader instead of deciding on their behalf.

**Goal:** filter a feed on the reader's own terms, using models that run on a post *before*
it is shown rather than after it is reported.

**Feature targets:**
- **Toxicity classification** - score every incoming post for hate, harassment, abuse. Blur above threshold before render.
- **NSFW tagging** - detect explicit or graphic images and text. Hide behind a warning the reader can lift.
- **Semantic trigger filter** - users describe topics to avoid in plain language; filter by embedding similarity, not keywords.
- **Pre-post toxicity check** - run the classifier on the draft in the compose box, offer a rephrase before it goes out.

**Delivery:** browser extension, mobile app, or web client.

**Hard constraint:** the majority of features must run on models *you trained yourself*,
served from the client or a backend you host. Commercial LLM APIs must not be doing all the work.

**Judges reward:** real trained models, not API wrappers. One well-tuned classifier working
on a live demo feed beats four half-built features calling GPT.

---

## Track 4 - Terminal Velocity: CLI Tools for Developer Workflows

Developers waste hours on repetitive terminal tasks: setting up projects, switching
environments, cleaning branches, digging through logs. Aliases and one-off scripts patch the
problem but rarely get shared or maintained. Build a CLI tool that removes a real, recurring
pain point. Fast, scriptable, something you would actually keep installed.

**All four deliverables are mandatory:**
1. **Working CLI tool** - installs with one command, completes at least one workflow end to end, supports `--help` with clear errors.
2. **Public repo** - README covering installation and usage examples.
3. **Demo** - max 3 minutes, showing the workflow before and after the tool.
4. **Design note** - max 1 page: the problem, key design choices, limitations.

**Judges reward:** a pain point the builder has personally felt. The demo should make a judge
think "I want this." Rough edges are fine. One thing done fast and reliably beats a Swiss
Army knife that jams.
