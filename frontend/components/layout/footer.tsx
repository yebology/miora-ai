export function Footer() {
  return (
    <footer className="border-t px-6 py-6">
      <div className="mx-auto flex max-w-5xl flex-col items-center gap-3 text-xs text-muted-foreground sm:flex-row sm:justify-between">
        <span>Miora AI</span>
        <span>Not financial advice. Trade at your own risk.</span>
        <div className="flex gap-4">
          <a
            href="https://github.com/yebology/miora-ai"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground"
          >
            GitHub
          </a>
          <a
            href="https://linkedin.com/in/yobelnathanielfilipus"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground"
          >
            LinkedIn
          </a>
        </div>
      </div>
    </footer>
  );
}
