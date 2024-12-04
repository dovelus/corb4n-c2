//import React from "react";
import { Link } from "react-router-dom";
import { Home, User, Menu } from "lucide-react";
import { Button } from "./ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Separator } from "./ui/separator";
import { ModeToggle } from "./mode-toggle";

interface TopbarProps {
  onMobileMenuToggle: () => void;
}

export function Topbar({ onMobileMenuToggle }: TopbarProps) {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-14 items-center justify-between">
          <div className="flex items-center space-x-4">
            <Link to="/" className="flex items-center space-x-2">
              <Home className="h-5 w-5" />
              <span className="sr-only">Home</span>
            </Link>
            <nav className="hidden md:flex items-center space-x-4">
              <Button asChild variant="ghost" size="sm">
                <Link to="/">Dashboard</Link>
              </Button>
              <Button asChild variant="ghost" size="sm">
                <Link to="/settings">Settings</Link>
              </Button>
              <Button asChild variant="ghost" size="sm">
                <Link to="/status">Status</Link>
              </Button>
            </nav>
          </div>
          <div className="flex items-center space-x-4">
            <ModeToggle />
            <Separator orientation="vertical" className="h-5" />
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  className="relative h-8 w-8 rounded-full"
                >
                  <User className="h-4 w-4" />
                  <span className="sr-only">User menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem>Logged in as: John Doe</DropdownMenuItem>
                <DropdownMenuItem>
                  <Link to="/logout">Log out</Link>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={onMobileMenuToggle}
            >
              <Menu className="h-5 w-5" />
              <span className="sr-only">Toggle menu</span>
            </Button>
          </div>
        </div>
      </div>
    </header>
  );
}
