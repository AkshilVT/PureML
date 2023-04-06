import clsx from "clsx";
import { File, Box, Database, ChevronDown, Menu, X } from "lucide-react";
import Button from "./ui/Button";
import { useMatches, useNavigate } from "@remix-run/react";
import { tv, type VariantProps } from "tailwind-variants";
import AvatarIcon from "./ui/Avatar";
import Dropdown from "./ui/Dropdown";
import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import { useState } from "react";

function linkCss(currentPage: boolean) {
  return clsx(
    currentPage ? " text-brand-200 " : " text-slate-500 ",
    " hover:text-brand-200 flex justify-center items-center px-5 cursor-pointer "
  );
}

const navbarStyles = tv({
  base: "navbar px-12 2xl:pr-0  max-w-screen-2xl",
  variants: {
    intent: {
      loggedIn: "",
      loggedOut: "",
    },
    fullWidth: {
      true: "w-full",
    },
  },
  defaultVariants: {
    intent: "loggedOut",
    fullWidth: true,
  },
});
interface Props extends VariantProps<typeof navbarStyles> {
  user: string;
  orgName: any;
  orgAvatarName: string;
}

export default function NavBar(props: Props) {
  const navigate = useNavigate();
  const matches = useMatches();
  const pathname = matches[1].pathname;
  const [open, setOpen] = useState(true);
  if (!open)
    return (
      <ul className="bg-white rounded-b-2xl border-b border-slate-200">
        <li className="flex py-4 px-12 md:px-12 w-full justify-between max-w-7xl">
          <li className="flex justify-center items-center pr-12">
            {props.intent === "loggedIn" && (
              <AvatarIcon>{props.orgAvatarName}</AvatarIcon>
            )}
            <div className="px-2 font-medium text-slate-600">
              {props.orgName}
            </div>
            {props.intent === "loggedIn" && (
              <DropdownMenu.Root>
                <DropdownMenu.Trigger className="flex justify-between items-center z-50 focus:outline-none z-50">
                  <ChevronDown className="text-slate-400 w-4" />
                </DropdownMenu.Trigger>
                <DropdownMenu.Content
                  sideOffset={7}
                  align="center"
                  className="bg-white justify-center items-center text-slate-600 rounded w-44 z-50 shadow focus:outline-none"
                >
                  <DropdownMenu.Item
                    className="flex px-3 py-3 justify-left items-center rounded outline-none hover:bg-slate-100 cursor-pointer focus:outline-none"
                    onClick={() => {
                      navigate(`/switchOrg`);
                    }}
                  >
                    Switch Organization
                  </DropdownMenu.Item>
                  {/* <DropdownMenu.Item
                      className="flex px-3 py-3 justify-left items-center rounded outline-none hover:bg-slate-100 cursor-pointer focus:outline-none"
                      onClick={() => {
                        navigate(`/joinOrg`);
                      }}
                    >
                      Join Organization
                    </DropdownMenu.Item> */}
                  <DropdownMenu.Item
                    className="flex px-3 py-3 justify-left items-center rounded outline-none hover:bg-slate-100 cursor-pointer focus:outline-none"
                    onClick={() => {
                      navigate(`/createOrg`);
                    }}
                  >
                    Create Organization
                  </DropdownMenu.Item>
                </DropdownMenu.Content>
              </DropdownMenu.Root>
            )}
          </li>
          <X
            className="sm:hidden text-slate-900 cursor-pointer w-8 h-8"
            onClick={() => setOpen(!open)}
          />
        </li>
        <li className="flex flex-col gap-y-2 px-12 py-4 font-medium text-brand-200">
          <li className="flex items-center">
            <a href="/models" className="w-max">
              Models
            </a>
          </li>
          <li className="flex items-center">
            <a href="/datasets" className="w-max">
              Datasets
            </a>
          </li>
          <li className="flex items-center">
            <a href="https://pureml.mintlify.app" className="w-max">
              Docs
            </a>
          </li>
          {props.intent === "loggedOut" ? (
            <li className="flex">
              <div className="w-full flex justify-center items-center px-5 font-medium text-slate-500">
                <a href="/auth/signin" className="w-max">
                  Sign in
                </a>
              </div>
              <div className="w-fit">
                <Button intent="primary">Sign up</Button>
              </div>
            </li>
          ) : (
            <li className="w-full flex items-center">
              <Dropdown intent="primary">
                <AvatarIcon>{props.user}</AvatarIcon>
              </Dropdown>
            </li>
          )}
        </li>
      </ul>
    );

  return (
    <div className="flex justify-center bg-slate-50 border-b border-slate-100">
      <div className={navbarStyles(props)}>
        <ul className="flex w-full justify-between">
          <li className="flex justify-center items-center pr-12">
            {props.intent === "loggedIn" && (
              <AvatarIcon>{props.orgAvatarName}</AvatarIcon>
            )}
            <div className="px-2 text-slate-600 font-medium">
              {props.orgName}
            </div>
            {props.intent === "loggedIn" && (
              <DropdownMenu.Root>
                <DropdownMenu.Trigger className="flex justify-between items-center z-50 focus:outline-none z-50">
                  <ChevronDown className="text-slate-400 w-4" />
                </DropdownMenu.Trigger>
                <DropdownMenu.Content
                  sideOffset={7}
                  align="center"
                  className="bg-white text-sm justify-center items-center text-xs text-slate-600 rounded w-44 z-50 shadow focus:outline-none"
                >
                  <DropdownMenu.Item
                    className="flex px-3 py-2 text-sm justify-left items-center rounded outline-none hover:bg-slate-100 cursor-pointer focus:outline-none"
                    onClick={() => {
                      navigate(`/switchOrg`);
                    }}
                  >
                    Switch Organization
                  </DropdownMenu.Item>
                  {/* <DropdownMenu.Item
                    className="flex px-3 py-2 text-sm justify-left items-center rounded outline-none hover:bg-slate-100 cursor-pointer focus:outline-none"
                    onClick={() => {
                      navigate(`/joinOrg`);
                    }}
                  >
                    Join Organization
                  </DropdownMenu.Item> */}
                  <DropdownMenu.Item
                    className="flex px-3 py-2 text-sm justify-left items-center rounded outline-none hover:bg-slate-100 cursor-pointer focus:outline-none"
                    onClick={() => {
                      navigate(`/createOrg`);
                    }}
                  >
                    Create Organization
                  </DropdownMenu.Item>
                </DropdownMenu.Content>
              </DropdownMenu.Root>
            )}
          </li>
          <li className="md:hidden">
            <Menu
              className="sm:hidden text-slate-900 cursor-pointer w-8 h-8"
              onClick={() => setOpen(!open)}
            />
          </li>
          {/* <Input
            intent="search"
          /> */}
        </ul>

        <ul className="hidden md:flex md:justify-center items-center">
          <li
            onClick={() => {
              navigate(`/models`);
            }}
            className={`${linkCss(pathname === `/models`)}`}
          >
            <Box className="w-4 h-4" />
            <span className="pl-2 font-medium">Models</span>
          </li>

          <li
            onClick={() => {
              navigate(`/datasets`);
            }}
            className={`${linkCss(pathname === `/datasets`)}`}
          >
            <Database className="w-4 h-4" />
            <span className="pl-2 font-medium">Datasets</span>
          </li>

          <a
            href="https://pureml.mintlify.app"
            className="flex justify-center items-center cursor-pointer px-5 hover:text-brand-200 border-r-2 border-slate-slate-200 font-medium text-slate-500"
          >
            <File className="w-4 h-4" />
            <span className="pl-2">Docs</span>
          </a>
          {props.intent === "loggedOut" ? (
            <li className="flex">
              <div className="w-full flex justify-center items-center px-5 font-medium text-slate-500">
                <a href="/auth/signin" className="w-max">
                  Sign in
                </a>
              </div>
              <div className="w-fit">
                <Button intent="primary">Sign up</Button>
              </div>
            </li>
          ) : (
            <li className="w-full flex justify-center items-center pl-5">
              <Dropdown intent="primary">
                <AvatarIcon>{props.user}</AvatarIcon>
              </Dropdown>
            </li>
          )}
        </ul>
      </div>
    </div>
  );
}
