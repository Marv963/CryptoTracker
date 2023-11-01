"use client";

import React from "react";
import {
  Dropdown,
  DropdownTrigger,
  DropdownMenu,
  DropdownItem,
  Button,
} from "@nextui-org/react";
import { MoonIcon, SunIcon } from "lucide-react";
import { useTheme } from "next-themes";
// import { FaGear } from "react-icons/fa6";

export default function ModeToggle() {
  const { setTheme } = useTheme();
  const iconClasses =
    "text-xl text-default-500 pointer-events-none flex-shrink-0";

  return (
    <Dropdown backdrop="blur">
      <DropdownTrigger>
        <Button isIconOnly aria-label="Like" color="secondary">
          <SunIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
          <MoonIcon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
        </Button>
      </DropdownTrigger>
      <DropdownMenu variant="faded" aria-label="Dropdown menu with icons">
        <DropdownItem
          key="light"
          startContent={<SunIcon className={iconClasses} />}
          onClick={() => setTheme("light")}
        >
          Light
        </DropdownItem>
        <DropdownItem
          key="dark"
          startContent={<MoonIcon className={iconClasses} />}
          onClick={() => setTheme("dark")}
        >
          Dark
        </DropdownItem>
        {/* <DropdownItem */}
        {/*   key="system" */}
        {/*   startContent={<FaGear className={iconClasses} />} */}
        {/*   onClick={() => setTheme("system")} */}
        {/* > */}
        {/*   System */}
        {/* </DropdownItem> */}
      </DropdownMenu>
    </Dropdown>
  );
}
