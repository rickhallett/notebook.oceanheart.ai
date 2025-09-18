---
title: "Building Passport: One Cookie to Rule Them All"
date: "2025-09-15"
tags: ["sso", "authentication", "security", "jwt", "go", "cognitive-skill:systems-thinking"]
summary: "A deep dive into Oceanheart Passport, the foundational identity layer for my entire ecosystem. Why I built a custom SSO and what it means to establish a perimeter of trust."
draft: false
---

# Building Passport: One Cookie to Rule Them All

Every ecosystem, digital or natural, needs a clear understanding of who's in it and what they're allowed to do. When I started mapping out the Oceanheart projects, I knew I couldn't build four separate apps with four separate logins. That's a recipe for friction and a security nightmare. The very first pillar had to be the foundation of trust for the whole world: **Oceanheart Passport**.

I call it the central bouncer for my domain. It's a Single Sign-On (SSO) service, but the philosophy runs deeper than just convenience. It's about creating a single, hardened perimeter where identity is managed with extreme care. Log in once, and a chic little cookie named `oh_session` becomes your verified pass to wander everywhere else—from my Notebook to the sensitive clinical tools.

## Identity Without Friction

The core user experience had to be seamless. The last thing I want is for someone to have a moment of insight in Sidekick, only to be jarringly interrupted by a login screen when they want to document it in Notebook. That's where the `oh_session` cookie comes in. It's a JWT (JSON Web Token), a kind of digital wristband that proves "this human's cool" to every other service in the `oceanheart.ai` domain.

But convenience can't come at the cost of security. One of the scariest vulnerabilities on the web is an "open redirect," where a malicious link can trick a login service into sending you to a phishing site after you sign in. To combat this, Passport has an aggressive validation system for its `returnTo` parameter. It ensures that after you log in, you are only ever sent back to a legitimate, intended destination within my ecosystem. No wandering off into evil portals.

The API is deliberately RESTful and emotionally mature: you `POST` to sign in, `DELETE` to sign out. It's clean, predictable, and does one job exceptionally well.

## The Why: An Architecture of Trust

Building a custom SSO isn't the easy route, but it was a non-negotiable architectural decision. By centralizing authentication, I can pour all my security efforts into one place. Instead of managing four separate, potentially weaker guard posts, I'm building a single fortress.

This is where the **Role-Based Access Control (RBAC)** in Passport's admin panel becomes critical. It's not just about *who* you are, but *what you can do*. A clinician needs full read/write access in Watson, but a researcher might only get access to anonymized, public data. Passport is the central authority that enforces these boundaries across the entire system.

By consolidating this logic, I can implement more sophisticated security measures—like advanced intrusion detection or multi-factor authentication—at a level that would be impractical to replicate across every single app. It's a design choice that accepts a single point of failure but works tirelessly to make that point as resilient as humanly possible. This is the bedrock. Everything else is built on the trust that Passport establishes.