import { Menu } from "@headlessui/react";
import { EllipsisVerticalIcon } from "@heroicons/react/24/solid"; // updated import for Heroicons v2

function ActionMenu({ actions = [] }) {
  return (
    <Menu as="div" className="relative inline-block text-left">
      <Menu.Button className="p-2 hover:bg-gray-100 rounded-full">
        <EllipsisVerticalIcon className="w-5 h-5 text-gray-700" />
      </Menu.Button>

      <Menu.Items className="absolute right-0 mt-2 w-40 bg-white border rounded-md shadow-lg z-[999]">
        {actions.map((action, index) => (
          <Menu.Item key={index}>
            {({ active }) => (
              <button
                onClick={action.onClick}
                className={`w-full text-left px-4 py-2 text-sm ${
                  active ? "bg-gray-100" : ""
                }`}
              >
                {action.label}
              </button>
            )}
          </Menu.Item>
        ))}
      </Menu.Items>
    </Menu>
  );
}

export default ActionMenu;
