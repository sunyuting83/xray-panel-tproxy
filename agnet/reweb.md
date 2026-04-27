Here is the professional, full-English instruction set for your React refactor agent. I’ve tailored it to be high-authority and precise, ensuring the AI stays within the `xpanel-web` boundary.

---

# 🚀 React 18.2.0 to Latest Evolution Agent Instructions

## 🎯 Core Objective
Upgrade the `xpanel-web` project from **React 18.2.0** to the **latest stable production version** (e.g., React 19). The migration must modernize syntax and performance APIs without altering existing business logic or API configurations.

---

## 🛡️ Strict Operational Constraints
1.  **Directory Jail**: 
    * **DO NOT** access any directories outside of `xpanel-web`.
    * **DO NOT** modify or enter `node_modules`.
    * **READ ACCESS** allowed for `package.json` only to verify dependencies.
2.  **Package Management**:
    * **MANDATORY** use of `yarn` for all installations, upgrades, and lockfile maintenance.
3.  **Code Integrity**:
    * **STRICTLY PROHIBITED**: Do not modify business logic within page components.
    * **STRICTLY PROHIBITED**: Do not touch API configuration files (e.g., `api.js`, `request.ts`, or any endpoints setup).
    * **STRICTLY PROHIBITED**: Do not alter UI layouts or existing CSS/SCSS modules.

---

## 🛠️ Refactoring Task List

### 1. Core Dependency Upgrade
* Execute: `yarn add react@latest react-dom@latest`.
* Synchronize Type Definitions: `yarn add -D @types/react@latest @types/react-dom@latest`.
* Audit `peerDependencies`: Identify conflicts with 3rd-party libraries (e.g., `antd`, `react-router-dom`). Report conflicts before attempting minor version bumps.

### 2. Root API Modernization
* Ensure the entry point (`index.js` or `main.tsx`) utilizes the current **Root API**:
    ```javascript
    import { createRoot } from 'react-dom/client';
    const container = document.getElementById('root');
    const root = createRoot(container!); 
    root.render(<App />);
    ```

### 3. Syntax & Hook Migration
* **Legacy Cleanup**: Remove deprecated lifecycle methods or legacy patterns (e.g., `string refs`, `findDOMNode`).
* **React 19 Readiness**: 
    * If upgrading to v19, simplify `forwardRef` usage as `ref` is now a standard prop.
    * Evaluate implementation of `useActionState` or `useOptimistic` only where they replace equivalent custom logic without changing behavior.
* **Concurrency Audit**: Wrap state-heavy updates in `useTransition` where appropriate to leverage the Concurrent Renderer.

### 4. Fragment & Cleanup
* Ensure `useEffect` cleanup functions are optimized for the latest strict mode behaviors.
* Convert appropriate `null` returns to `<Fragment />` or `<></>` if it improves component consistency.

---

## 📋 Execution Protocol
1.  **Discovery**: Read `package.json` to map the current tech stack.
2.  **Pruning**: Clean up `yarn.lock` for redundant or conflicting entries.
3.  **Installation**: Run `yarn` commands to update core packages.
4.  **Transformation**: Scan `.js/.jsx/.ts/.tsx` files in `/src`. Apply syntactic sugar and modern Hook patterns.
5.  **Verification**: Confirm that all refactored components still reference the original API/service files and that paths remain intact.

---

## 🚩 Error Handling & Escalation
* If a 3rd-party library (e.g., a specific UI component lib) is fundamentally incompatible with the latest React version, the Agent **must halt** and request human intervention. **Do not** use `--force` or `--ignore-engines` without explicit permission.

---

**Please acknowledge these instructions. If understood, begin by reading `xpanel-web/package.json`.**