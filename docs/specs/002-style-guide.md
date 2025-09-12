# Hugo Mini Theme Style Guide

**Document**: Notebook.OceanHeart.ai Style Guide  
**Version**: 1.0  
**Date**: 2025-09-12  
**Based on**: hugo-theme-mini by @nodejh

## Overview

This style guide reverse-engineers the Hugo Mini theme into a framework-agnostic design system for the notebook.oceanheart.ai blog engine. The guide focuses on visual design, typography, layout principles, and component styling while excluding Hugo-specific templating and i18n features.

## Design Philosophy

Hugo Mini follows a **minimalist aesthetic** with emphasis on:
- Clean typography and generous whitespace
- Subtle color palette with blue accents
- Responsive design with mobile-first approach
- Content-focused layout with minimal distractions
- Accessible and semantic HTML structure

## Typography

### Base Typography
```css
body {
  font: normal 15px/1.5em 'Helvetica Neue', Helvetica, Arial, sans-serif;
  color: #404040;
  line-height: 1.75;
  letter-spacing: 0.008em;
}
```

### Headings
- **Font Weight**: 400 (regular)
- **Color**: #404040 (same as body text)
- **Philosophy**: Headings differentiated by size, not weight or color

### Code Typography
```css
/* Inline Code */
code {
  font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, Courier, monospace;
  background-color: rgba(0, 0, 0, 0.06);
  padding: 0 2px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 2px;
}

/* Code Blocks */
pre code {
  font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, Courier, monospace;
}
```

## Color Palette

### Primary Colors
- **Text**: #404040 (dark gray)
- **Secondary Text**: #757575 (medium gray)
- **Light Text**: #bbb, #8c8c8c
- **Background**: #fff (white)

### Accent Colors
- **Link Blue**: #5badf0 (primary)
- **Link Hover**: #0366d6 (darker blue)
- **Button Background**: #5badf0

### UI Colors
- **Borders**: #eee (light), #dadadc (medium)
- **Table Alternate**: #fcfcfc
- **Code Background**: rgba(0, 0, 0, 0.06)
- **Blockquote**: #eee (border), #757575 (text)

## Layout System

### Container Widths
```css
/* Different content areas have specific max-widths */
#list-page { max-width: 580px; }    /* Home/List pages */
#single { max-width: 680px; }       /* Article pages */
#archive { max-width: 580px; }      /* Archive pages */
#tags { max-width: 700px; }         /* Tag pages */
```

### Spacing Scale
- **Small**: 12px, 18px, 20px
- **Medium**: 24px, 36px, 40px  
- **Large**: 48px, 60px, 64px
- **XLarge**: 100px

### Responsive Breakpoints
- **Mobile**: max-width: 700px
- **Small Mobile**: max-width: 324px

## Component Specifications

### 1. Site Header / Profile
```css
.profile {
  margin: 60px auto 0 auto;
  text-align: center;
}

.profile .avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
}

.profile h1 {
  font-size: 20px;
  color: #404040;
  margin-bottom: 0;
  margin-top: 0;
}

.profile h2 {
  font-size: 20px;
  font-weight: 300;
  color: #757575;
  margin-top: 0;
}
```

### 2. Navigation
```css
nav.navigation {
  padding: 20px 20px 0;
  background: rgba(255, 255, 255, 0.9);
  text-align: right;
}

nav.navigation a {
  padding: 8px 12px;
  color: #5badf0;
  font-size: 13px;
  border-radius: 3px;
}

nav.navigation a.button {
  background: #5badf0;
  color: #fff;
  margin-left: 12px;
}
```

### 3. Post List (Home Page)
```css
#list-page .item {
  margin: 12px 0;
}

#list-page .title {
  color: #404040;
  font-size: 20px;
  font-weight: 400;
  width: 80%;
}

#list-page .date {
  width: 20%;
  float: right;
  text-align: right;
  color: #bbb;
}

#list-page .summary {
  color: #757575;
  margin-top: 12px;
  margin-bottom: 36px;
}
```

### 4. Single Post Layout
```css
#single {
  margin: 60px auto 0 auto;
  padding: 0 64px;
}

#single .title {
  text-align: center;
  font-size: 32px;
  font-weight: 400;
  line-height: 48px;
}

#single .tip {
  text-align: center;
  color: #8c8c8c;
  margin-top: 18px;
  font-size: 14px;
}

#single .content {
  margin-top: 36px;
}
```

### 5. Table of Contents
```css
.toc {
  background: #f8f8f8;
  padding: 10px 7px;
  margin-top: 36px;
}

.toc details summary {
  cursor: zoom-in;
  margin-inline-start: 14px;
  font-weight: 500;
}

.toc details[open] summary {
  cursor: zoom-out;
}
```

### 6. Tags
```css
.tags a {
  background: #f2f2f2;
  padding: 4px 7px;
  color: #757575;
  font-size: 14px;
  margin-right: 3px;
}

.tags a:hover {
  color: #0366d6;
}
```

### 7. Archive/List Layout
```css
#archive .group {
  margin: 24px auto;
}

#archive .group .key {
  font-size: 20px;
  margin-bottom: 12px;
}

#archive .group .value .date {
  color: #999;
  min-width: 60px;
}

#archive .group .value .title a {
  color: #404040;
}
```

### 8. Pagination
```css
#list-page .pagination {
  margin: 48px 0;
  width: 100%;
  height: 32px;
}

#list-page .pagination .pre {
  float: left;
}

#list-page .pagination .next {
  float: right;
}
```

## Markdown Content Styling

### Blockquotes
```css
blockquote {
  margin-top: 5px;
  margin-bottom: 5px;
  padding-left: 1em;
  border-left: 3px solid #eee;
  color: #757575;
}
```

### Horizontal Rules
```css
hr {
  border: none;
  height: 2px;
  margin: 40px auto;
  background: #eee;
}
```

### Tables
```css
table {
  width: 100%;
  margin: 40px 0;
  border-collapse: collapse;
  line-height: 1.5em;
}

th {
  padding: 10px 15px;
  font-size: 13px;
  font-weight: bold;
  color: #444;
  border: 1px solid #dadadc;
}

td {
  padding: 10px 15px;
  border: 1px solid #dadadc;
}

tr:nth-child(odd) > td {
  background: #fcfcfc;
}
```

### Images
```css
.markdown-image img {
  max-width: 100%;
}
```

### Anchor Links
```css
.anchor { 
  font-size: 100%; 
  visibility: hidden; 
  color: silver;
}

h1:hover a,
h2:hover a,
h3:hover a,
h4:hover a { 
  visibility: visible;
}
```

## Dark Mode Support

Hugo Mini includes a simple dark mode implementation:
```css
html[theme='dark-mode'] {
  filter: invert(1) hue-rotate(180deg);
}
```

## Responsive Design Patterns

### Mobile Navigation (< 700px)
```css
@media (max-width: 700px) {
  nav.navigation {
    padding: 20px 10px 0 0;
  }
  
  nav.navigation a {
    padding: 8px 8px;
  }
}
```

### Mobile Single Post (< 700px)
```css
@media (max-width: 700px) {
  #single {
    padding: 0 18px;
    margin: 20px auto 0 auto;
  }
  
  #single .title {
    font-size: 24px;
    line-height: 32px;
  }
}
```

### Very Small Screens (< 324px)
```css
@media (max-width: 324px) {
  nav.navigation a.button {
    display: none;
  }
}
```

## Footer
```css
#footer {
  margin-top: 100px;
  margin-bottom: 100px;
  text-align: center;
  color: #bbbbbb;
  font-size: 14px;
}

#footer .copyright {
  margin: 20px auto;
  font-size: 15px;
}
```

## Implementation Guidelines

### CSS Organization
1. **Base styles**: Typography, colors, resets
2. **Layout**: Container widths, spacing, grid
3. **Components**: Navigation, profile, post lists, etc.
4. **Content**: Markdown styling, code blocks, tables
5. **Responsive**: Mobile-first media queries

### Key Principles
1. **Consistency**: Use the defined color palette and spacing scale
2. **Readability**: Maintain generous line-height and spacing
3. **Hierarchy**: Use size and spacing for visual hierarchy, not color/weight
4. **Performance**: Minimize CSS, use system fonts
5. **Accessibility**: Ensure sufficient contrast, semantic HTML

### Adaptation Notes for Notebook Engine
1. Replace Hugo templating with Go template syntax
2. Remove i18n-specific features
3. Adapt container widths for content strategy
4. Integrate with existing Chroma syntax highlighting
5. Ensure compatibility with current middleware stack

## Content Structure Requirements

### Post Metadata
- **Title**: Displayed prominently at top of post
- **Date**: Formatted and displayed in post tip section
- **Tags**: Displayed as clickable badges below content
- **Summary**: Used in post lists and meta descriptions

### Navigation Structure
- **Home**: Return to main page
- **Archive**: Chronological post listing
- **Tags**: Tag cloud/index
- **About**: Static page
- **RSS**: Feed subscription (optional button)

This style guide provides the foundation for implementing the Hugo Mini aesthetic in the notebook.oceanheart.ai blog engine while maintaining the clean, minimalist design philosophy that makes the theme effective.