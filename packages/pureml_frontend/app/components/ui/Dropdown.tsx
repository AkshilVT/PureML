import { tv, type VariantProps } from "tailwind-variants";
import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import { ChevronDown } from "lucide-react";
import { useNavigate } from "@remix-run/react";
import { toast } from "react-toastify";
import type { ReactNode } from "react";

const dropdownStyles = tv({
  slots: {
    base: "focus:outline-none z-50",
    content:
      "bg-white justify-center items-center text-sm text-slate-600 rounded w-44 z-50 shadow",
    items:
      "flex px-3 py-2 text-sm justify-left items-center rounded outline-none hover:bg-slate-100 cursor-pointer",
  },
  variants: {
    intent: {
      primary: "flex justify-between items-center",
      orgType:
        "flex justify-between items-center text-sm text-slate-600 rounded border border-slate-600 h-8 px-2",
    },
    fullWidth: {
      true: "w-full",
      false: "w-28",
    },
  },
  defaultVariants: {
    intent: "primary",
    fullWidth: true,
  },
});

export interface Props extends VariantProps<typeof dropdownStyles> {
  [children: string]: ReactNode;
}

export default function Dropdown(props: Props) {
  const navigate = useNavigate();
  const { base, content, items } = dropdownStyles();

  return (
    <div>
      <DropdownMenu.Root>
        <DropdownMenu.Trigger className={base()}>
          {props.children}
          {props.intent !== "primary" ? (
            <ChevronDown className="text-slate-400" />
          ) : (
            ""
          )}
        </DropdownMenu.Trigger>
        {props.intent === "primary" ? (
          <DropdownMenu.Content
            sideOffset={7}
            align="end"
            className={content()}
          >
            <DropdownMenu.Item
              className={items()}
              onClick={() => {
                navigate(`/profile`);
              }}
            >
              Profile
            </DropdownMenu.Item>
            <DropdownMenu.Separator className="h-px bg-slate-200" />
            <DropdownMenu.Item
              className={items()}
              onClick={() => {
                navigate(`/settings`);
              }}
            >
              Settings
            </DropdownMenu.Item>
            {/* <DropdownMenu.Item
              className={items()}
              onClick={() => {
                navigate(`/contact`);
              }}
            >
              Contact Us
            </DropdownMenu.Item> */}
            <DropdownMenu.Separator className="h-px bg-slate-200" />
            <DropdownMenu.Item
              className={items()}
              onClick={() => {
                toast.success("Sad to see you go :(");
                navigate(`/logout`);
              }}
            >
              Sign out
            </DropdownMenu.Item>
          </DropdownMenu.Content>
        ) : (
          <DropdownMenu.Content sideOffset={7} className={content()}>
            <DropdownMenu.Item className={items()}>Company</DropdownMenu.Item>
            <DropdownMenu.Item className={items()}>
              University
            </DropdownMenu.Item>
            <DropdownMenu.Item className={items()}>Classroom</DropdownMenu.Item>
            <DropdownMenu.Item className={items()}>
              Non-profit
            </DropdownMenu.Item>
            <DropdownMenu.Item className={items()}>Community</DropdownMenu.Item>
          </DropdownMenu.Content>
        )}
      </DropdownMenu.Root>
    </div>
  );
}
