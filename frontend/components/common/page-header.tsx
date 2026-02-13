import { cn } from "@/lib/utils";

interface PageHeaderProps {
  title: string;
  description?: string;
  action?: React.ReactNode;
  className?: string;
  children?: React.ReactNode;
}

export function PageHeader({
  title,
  description,
  action,
  className,
  children,
}: PageHeaderProps) {
  return (
    <div
      className={cn(
        "flex flex-col gap-3 sm:gap-4 md:flex-row md:items-center md:justify-between",
        className,
      )}
    >
      <div className="space-y-0.5 sm:space-y-1 min-w-0">
        <h1 className="text-xl sm:text-2xl font-bold tracking-tight md:text-3xl truncate">
          {title}
        </h1>
        {description && (
          <p className="text-sm sm:text-base text-muted-foreground line-clamp-2">
            {description}
          </p>
        )}
      </div>
      {action && (
        <div className="flex items-center gap-2 flex-shrink-0">{action}</div>
      )}
      {children}
    </div>
  );
}
