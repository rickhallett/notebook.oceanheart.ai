---
title: "Anatomy of a Wild Day: Debugging, Docker, and Digital Archaeology"
date: "2025-09-13"
tags: ["debugging", "hsts", "full-stack", "devops", "cognitive-skill:debugging", "cognitive-skill:context-switching", "cognitive-skill:systems-thinking"]
summary: "A play-by-play of one of those days in software development where a single bug leads to a deep dive into browser internals, followed by a whirlwind tour through half a dozen technologies—all before lunch."
draft: false
---

# One of *Those* Days

You know the ones. The day starts with a simple, routine task. For me, it was just trying to access my local Django admin page. Click, log in, done. Should've taken ten seconds.

Instead, I hit a wall. A digital dead end. The page was just unreachable.

What followed was a journey so deep into the guts of my web browser that it felt like digital archaeology. It was one of those developer mysteries that starts with a simple "why isn't this working?" and ends with you questioning the very fabric of reality.

## The Ghost in the Machine

My first instinct, like any developer, was to go through the standard checklist.

1.  Clear cookies. Nope.
2.  Nuke the cache. Nothing.
3.  Delete all history. Still broken.

I did everything short of degaussing my monitor. A friend later joked that I had "metaphorically cleared my chakras" trying to solve it, which is painfully accurate. I was desperate. The page just wouldn't load. It was like a door slamming shut every time I tried to open it.

The problem felt personal. The browser remembered something, a promise it made to a past version of my code, and it refused to let go. This wasn't just a bug; it was a stubborn, invisible rule that had decided how my world worked now.

## Meet the Villain: HSTS Pinning

After hours of digging through obscure forums and internal debug menus—places most people never see, like `chrome://net-internals`—I finally found the culprit. It wasn't a random bug. It was a security feature called **HSTS**, or HTTP Strict Transport Security.

HSTS's job is simple: force a browser to *only* use a secure HTTPS connection for a specific website. It protects against attacks where someone tries to downgrade you to an insecure connection. It’s a good thing!

But here's the twist. My local development server wasn't set up with HTTPS. At some point, probably through an accidental click on a link in a redirect chain, my browser had received two commands:

1.  **An HSTS Policy:** "For this address, *always* use HTTPS. No exceptions."
2.  **A 301 Redirect:** "The page you want has permanently moved to this *HTTPS* address."

My browser, trying to be both efficient and secure, combined these two rules into an unbreakable vow. It cached the redirect and locked it in with the HSTS policy. It decided that my local admin page was now HTTPS-only, forever. It wouldn't even *try* to load the HTTP version anymore. It just failed instantly.

The solution? Scorched earth. I had to navigate into hidden system folders (`~/Library/Application Support/Google/Chrome/Default` for anyone curious) and manually delete the `TransportSecurity` file—Chrome's little black book of HSTS rules. It felt less like typing and more like performing an exorcism.

And it worked. The page loaded. The relief was immense. I felt like I'd conquered a mountain.

## The Punchline: That Was Just My Morning Warm-Up

Here’s the crazy part. That whole debugging saga—the frustration, the deep dive, the triumphant fix—was all before 11 AM.

Looking at my Git history from that day is humbling. The HSTS fight was just the opening act. The rest of the morning was a high-speed whirlwind of building, documenting, and deploying.

- **10:27 AM:** Committing JWT cookie helpers. The foundation of a secure authentication system.
- **10:38 AM:** Adding the core JWT authentication logic.
- **10:55 AM:** Building a full-blown Role-Based Access Control (RBAC) system for the admin panel.

In less than 30 minutes, I’d built what a friend jokingly called "Role-Based Anxiety Control." An entire security and identity system, from scratch, before my first coffee break was even cold.

## A Tour of the Multiverse

The real challenge wasn't just the features, but the constant context-switching. My brain was jumping between completely different technology worlds, minute by minute.

- **The Python World:** I was deep in Django and FastAPI, setting up isolated environments (`venv`), configuring the Gunicorn web server, and using Caddy to handle traffic and automate security.
- **The Go World:** I switched over to Go to build file-based web templates, set up hot-reloading scripts for instant feedback (a lifesaver for development speed), and even built a secure admin reload endpoint.
- **The Frontend World:** A quick cameo from SvelteKit and Tailwind CSS to spin up a quick UI prototype, then back to the backend grind.

It felt less like coding and more like being a chef with four ovens going, cooking three different cuisines, with a single smoke alarm named Docker watching over everything.

## From Local Code to Global Cloud

And it wasn't just happening on my machine. I was building the entire deployment pipeline at the same time.

I was wrestling with multi-stage Dockerfiles, setting up CI workflows for automated database migrations, and configuring deployments to **Fly.io** (which spins up tiny servers close to users) and a distributed database with **Turso**. I wasn’t just building an app; I was building the factory that builds and runs the app, globally.

This is the new reality. You're not just a developer anymore. You're the architect, the operations engineer, and the quality control, all at once. The cognitive load is immense. It's a constant balancing act between building features, managing infrastructure, and documenting your decisions in PRDs and ADRs so you don't forget *why* you made a choice by 4 PM.

## The Takeaway: It's All Connected

That day was a perfect storm of modern software development. It was a reminder that persistence in browsers is a powerful, sometimes dangerous, force. It showed that the lines between technologies are blurring, and fluency across a wide stack is becoming the norm.

Most importantly, it highlighted that the biggest bottleneck isn't the code or the cloud—it's your own brain trying to keep the entire complex, interconnected system in your head at once.

And sometimes, even in the middle of that storm, you find a moment to do something small, like minifying a `package.json` file. Not because it makes the code faster, but for the simple, human need for order and craftsmanship. Even when you're building a digital universe, the little details still matter.