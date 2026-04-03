# Agentic Coding Guidelines

This document provides comprehensive guidelines for AI coding agents working on this React + TypeScript + Vite frontend application. Follow these conventions to maintain consistency and quality across the codebase.

## Project Overview

- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite 8
- **Styling**: TailwindCSS with shadcn/ui components
- **Routing**: React Router v7
- **Form Handling**: react-hook-form with zod validation
- **State Management**: React hooks (useState, useEffect, useContext)
- **HTTP Client**: Fetch API with custom error handling

## Build & Development Commands

### Local Development
```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build locally
npm run preview

# Run type checking
npm run type-check

# Run linting
npm run lint
```

### Environment Variables
The application uses environment variables prefixed with `VITE_`:
- `VITE_API_BASE_URL` - Base URL for API endpoints
- All VITE_* variables are exposed to the client-side code

## Testing Guidelines

### Testing Framework
Currently no test framework is configured. When adding tests:
- Use Vitest as the testing framework (aligned with Vite)
- Place test files alongside source files with `.test.ts` or `.test.tsx` extension
- Use React Testing Library for component tests
- Mock API calls using `vi.mock()` or MSW (Mock Service Worker)

### Test Structure
```
src/
├── components/
│   ├── Button.tsx
│   └── Button.test.tsx
├── pages/
│   ├── LoginPage.tsx
│   └── LoginPage.test.tsx
└── utils/
    ├── api.ts
    └── api.test.ts
```

### Key Testing Scenarios
- Form validation with zod schemas
- API error handling and loading states
- Route navigation and protected routes
- Component interactions with shadcn/ui

## Code Style Guidelines

### TypeScript Configuration
- Strict mode enabled (`"strict": true`)
- Target: ES2023
- Module resolution: Node with path aliases
- Path alias: `@/*` maps to `src/*`
- No implicit any types allowed
- All functions must have explicit return types when not obvious

### React Patterns

#### Component Structure
- Use functional components with TypeScript interfaces/types
- Define props using `interface` (not `type`) for better error messages
- Keep components focused and single-responsibility
- Use named exports only (no default exports)

```tsx
// Good
interface ButtonProps {
  variant?: 'primary' | 'secondary';
  onClick: () => void;
  children: React.ReactNode;
}

export function Button({ variant = 'primary', onClick, children }: ButtonProps) {
  // implementation
}
```

#### Hooks Usage
- Custom hooks should start with `use` prefix
- Handle loading and error states explicitly
- Use `useCallback` and `useMemo` judiciously for performance
- Clean up effects properly in `useEffect`

#### State Management
- Prefer local component state (`useState`) for UI state
- Use context for global state that affects multiple components
- Avoid prop drilling by creating appropriate context providers
- Keep state as minimal as possible

### Routing (React Router v7)
- Define routes in `App.tsx` using `createBrowserRouter`
- Use `Link` component for navigation (not anchor tags)
- Handle route parameters with `useParams`
- Implement protected routes using wrapper components
- Use `useNavigate` for programmatic navigation

### Forms (react-hook-form + zod)
- Always define zod schemas for form validation
- Use `z.object()` with appropriate field validations
- Leverage `zodResolver` from `@hookform/resolvers/zod`
- Handle form submission with proper error boundaries
- Use shadcn/ui form components (`FormField`, `FormItem`, etc.)

```tsx
const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
});

const form = useForm<LoginFormValues>({
  resolver: zodResolver(loginSchema),
  defaultValues: { email: '', password: '' },
});
```

### API Integration
- Create dedicated API utility functions in `src/utils/api.ts`
- Handle HTTP errors with custom error types
- Use async/await pattern consistently
- Implement proper loading and error states in components
- Set appropriate headers (Content-Type, Authorization when needed)

```ts
interface ApiError extends Error {
  status: number;
  message: string;
}

export async function apiCall<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`API call failed: ${response.status} ${response.statusText}`);
  }

  return response.json();
}
```

## Styling Guidelines

### TailwindCSS
- Use semantic class names over presentational ones
- Leverage responsive prefixes (`sm:`, `md:`, `lg:`) appropriately
- Use Tailwind's spacing scale consistently (2, 4, 6, 8, etc.)
- Avoid arbitrary values when possible; extend theme if needed

### shadcn/ui Components
- Import components directly from their source paths:
  ```tsx
  import { Button } from '@/components/ui/button';
  import { Input } from '@/components/ui/input';
  ```
- Customize components using the `variant` and `size` props when available
- Extend components by wrapping them rather than modifying source
- Follow shadcn/ui accessibility patterns (proper ARIA attributes, keyboard navigation)

### Custom Utilities
- Use the provided `cn` utility function for conditional class merging:
  ```tsx
  import { cn } from '@/lib/utils';
  
  className={cn('base-class', condition && 'conditional-class')}
  ```
- The `cn` function combines `clsx` and `tailwind-merge` for safe class composition

## File Structure

```
src/
├── assets/           # Static assets (images, fonts)
├── components/       # Reusable UI components
│   └── ui/          # shadcn/ui components
├── hooks/           # Custom React hooks
├── lib/             # Utility functions and helpers
├── pages/           # Page-level components
├── routes/          # Route definitions (if separate)
├── styles/          # Global CSS and Tailwind configuration
├── types/           # Global TypeScript types and interfaces
└── utils/           # Utility functions (API, formatting, etc.)
```

## Naming Conventions

### Files
- Component files: `PascalCase.tsx` (e.g., `UserProfile.tsx`)
- Utility files: `camelCase.ts` (e.g., `formatDate.ts`)
- Hook files: `useCamelCase.ts` (e.g., `useLocalStorage.ts`)
- Type files: `PascalCase.types.ts` (e.g., `User.types.ts`)

### Variables and Functions
- Variables: `camelCase`
- Functions: `camelCase`
- Constants: `UPPER_SNAKE_CASE` (for truly constant values)
- React components: `PascalCase`
- Interfaces/Types: `PascalCase`

### Imports
- Use absolute imports with `@/` alias:
  ```tsx
  // Good
  import { Button } from '@/components/ui/button';
  import { apiCall } from '@/utils/api';
  
  // Avoid
  import { Button } from '../../components/ui/button';
  ```

## Error Handling

### Frontend Errors
- Display user-friendly error messages
- Log technical details to console for debugging
- Handle network errors gracefully
- Provide retry mechanisms where appropriate
- Use error boundaries for component-level error isolation

### Validation Errors
- Show inline validation errors next to form fields
- Use zod for schema validation
- Aggregate multiple errors appropriately
- Clear errors when user corrects input

## Performance Considerations

- Memoize expensive computations with `useMemo`
- Memoize callback functions with `useCallback` when passed to child components
- Implement virtualization for long lists
- Lazy load components with `React.lazy` and `Suspense`
- Optimize images and assets appropriately
- Avoid unnecessary re-renders with proper dependency arrays

## Accessibility (a11y)

- Use semantic HTML elements
- Provide proper ARIA labels and roles
- Ensure keyboard navigation works
- Maintain sufficient color contrast
- Support focus management
- Test with screen readers when possible

## Git Workflow

### Commit Messages
- Use conventional commits format: `type(scope): description`
- Common types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Keep commits small and focused
- Reference issues when applicable

### Branching
- Feature branches should follow `feature/descriptive-name` pattern
- Bug fixes should follow `fix/descriptive-name` pattern
- Always pull latest changes before creating PRs

## Security Considerations

- Never expose sensitive information in client-side code
- Sanitize user inputs appropriately
- Use HTTPS for all API calls
- Implement proper authentication flows
- Validate all data from external sources
- Keep dependencies updated

## Dependencies

### Adding New Dependencies
- Prefer lightweight, well-maintained packages
- Check bundle size impact before adding
- Ensure TypeScript support or provide types
- Verify license compatibility
- Update `package.json` and lock file consistently

### Key Dependencies to Know
- `react-hook-form`: Form state management
- `zod`: Schema validation
- `@radix-ui/react-*`: Accessible UI primitives
- `tailwind-merge`: Safe Tailwind class merging
- `clsx`: Conditional class name composition
- `class-variance-authority`: Variant-based class composition

## IDE Configuration

No specific Cursor rules or GitHub Copilot instructions exist in this repository. Agents should rely on the guidelines provided above and the existing codebase patterns.

When working with this codebase:
- Enable ESLint integration in your editor
- Use Prettier for consistent formatting
- Leverage TypeScript language service for type safety
- Follow the existing patterns demonstrated in the source files

## Getting Help

If you encounter ambiguous requirements or edge cases:
1. Look for similar patterns in existing code
2. Follow the principle of least surprise
3. Prioritize consistency with existing code over personal preferences
4. When in doubt, ask for clarification rather than making assumptions

Remember: The goal is to write code that fits seamlessly into the existing codebase while maintaining high quality, performance, and maintainability standards.

所有变更都需要集成测试
所有变更都需要通过集成测试
主要工作流程需要自动化测试
