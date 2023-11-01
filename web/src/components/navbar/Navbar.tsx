"use client";
import React from "react";
import {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
  Link,
  Button,
  NavbarMenuToggle,
  NavbarMenu,
  NavbarMenuItem,
} from "@nextui-org/react";
import Image from "next/image";
import { FaGithub } from "react-icons/fa6";
import ModeToggle from "../ui/mode-toggle";
import { usePathname } from "next/navigation";

export default function App() {
  const [isMenuOpen, setIsMenuOpen] = React.useState(false);
  const pathName = usePathname();

  const menuItems = [
    "Profile",
    "Dashboard",
    "Activity",
    "Analytics",
    "System",
    "Deployments",
    "My Settings",
    "Team Settings",
    "Help & Feedback",
    "Log Out",
  ];

  return (
    <Navbar onMenuOpenChange={setIsMenuOpen} isBordered>
      <NavbarContent>
        <NavbarMenuToggle
          aria-label={isMenuOpen ? "Close menu" : "Open menu"}
          className="sm:hidden"
        />
        <NavbarBrand>
          <div className="flex-shrink-0">
            <Link href="/">
              <div className="dark:hidden">
                <Image
                  src="/images/logo/light.svg"
                  height={28}
                  width={150}
                  alt="CryptoTracker"
                  className="h-auto"
                />
              </div>
              <div className="hidden dark:block">
                <Image
                  src="/images/logo/dark.svg"
                  height={28}
                  width={150}
                  alt="CryptoTracker"
                  className="h-auto"
                />
              </div>
            </Link>
          </div>
        </NavbarBrand>
      </NavbarContent>

      <NavbarContent className="hidden sm:flex gap-4" justify="center">
        <NavbarItem isActive={pathName === "/"}>
          <Link color="foreground" href="/">
            Coins
          </Link>
        </NavbarItem>
        <NavbarItem isActive={pathName === "/arbitrage"}>
          <Link href="/arbitrage" aria-current="page">
            Arbitrage
          </Link>
        </NavbarItem>
      </NavbarContent>
      <NavbarContent justify="end">
        <NavbarItem>
          <Button isIconOnly color="secondary" aria-label="Like">
            <FaGithub className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all" />
          </Button>
        </NavbarItem>
        <NavbarItem>
          <ModeToggle />
        </NavbarItem>
      </NavbarContent>
      <NavbarMenu>
        {menuItems.map((item, index) => (
          <NavbarMenuItem key={`${item}-${index}`}>
            <Link
              color={
                index === 2
                  ? "primary"
                  : index === menuItems.length - 1
                  ? "danger"
                  : "foreground"
              }
              className="w-full"
              href="#"
              size="lg"
            >
              {item}
            </Link>
          </NavbarMenuItem>
        ))}
      </NavbarMenu>
    </Navbar>
  );
}
