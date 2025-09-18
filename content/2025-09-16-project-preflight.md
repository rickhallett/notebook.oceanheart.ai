---
title: "Project Preflight: Readiness as a Ritual"
date: "2025-09-16"
tags: ["ai", "llm", "onboarding", "sveltekit", "fastapi", "cognitive-skill:user-empathy"]
summary: "Exploring Project Preflight, an AI-readiness tool that's less about checking boxes and more about preparing for a partnership. Why slowing down the first interaction is key to building trust."
draft: false
---

# Project Preflight: Readiness as a Ritual

Onboarding is one of the most neglected parts of software. We usually treat it as a race to get the user clicking buttons, a series of tooltips and "Got it!" prompts. But when the "tool" is a complex AI, especially one designed for sensitive contexts like mental health, that approach feels reckless. This conviction led me to build **Project Preflight**.

The goal of Preflight isn't to just *onboard* a user; it's to prepare them for a *partnership* with AI. The Readme describes it as an "AI-readiness questionnaire," but it's really a slow, reflective, and conversational process. It's designed to warm users into a relationship with the technology, setting expectations and understanding their comfort levels before they ever start a real session.

## Every Interaction is Research

Preflight is built on a simple but powerful mantra: **every interaction is research, and every survey is a conversation.**

Instead of a static form, the questions are driven by a versioned JSON file. This "Form DSL" means I can run experiments, change the flow, and adapt the survey without a full redeployment. The system autosaves progress and logs every interaction, turning the onboarding process itself into a valuable dataset.

After the initial questions, a carefully scripted LLM-powered chat begins. The prompt pipeline has firm rules: ask one question at a time, don't offer medical advice, stay polite. It's a gentle handshake, not an attempt to be a therapist. It gets the user accustomed to the conversational nature of the AI in a low-stakes environment.

## The Conviction: Slowing Down to Build Trust

In a world that prizes instant gratification, deliberately slowing a user down feels counterintuitive. I'm absolutely risking losing people who just want to get to the "point." But that's the core philosophical stance of Preflight.

For a tool as potentially impactful as Sidekick, a rushed, unconsidered engagement with its AI could be more harmful than helpful. Misunderstandings, frustration, and unrealistic expectations are real dangers. So, this "readiness as a ritual" is worth the potential friction because it does two things:

1.  **It prepares the user for the AI:** It encourages them to think introspectively about what they want from the interaction and what their boundaries are.
2.  **It prepares the AI for the user:** It gives the Oceanheart ecosystem a real-time, evolving profile of the user's comfort levels, biases, and learning style. This allows for an unprecedented level of personalized, ethical guardrails.

The tech stack is a diplomatic summit—a **SvelteKit** frontend for a great UX, a **FastAPI** backend for the logic, and **Passport** for auth. But the real technology here is the methodology: using the first point of contact not for extraction, but for mutual understanding. It's the first step in treating the human-AI interaction with the gravity it deserves.