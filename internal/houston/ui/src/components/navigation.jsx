import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs"
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar"
import { ChevronRight } from "lucide-react"
import Queues from "@/components/queues"

export default function Navigation({ activeTab = "queues", hideContent = false, breadcrumb = null }) {
  const handleTabChange = (value) => {
    // Use regular browser navigation
    const paths = {
      "queues": "/",
      "pubsub": "/pubsub",
      "users&access": "/users",
      "settings": "/settings"
    };
    
    if (paths[value]) {
      window.location.href = paths[value];
    }
  };

  return (
    <div className="">
      <div className="flex flex-row justify-between items-center px-6">
        <div>
          <a href="/" className="text-lg text-gray-500 hover:text-gray-900">PlainQ</a>
        </div>
        <div className="flex flex-row items-center gap-2">
          <div>
            <a href="https://docs.plainq.com">Docs</a>
          </div>
          <Avatar>
            <AvatarImage src="https://github.com/shadcn.png"/>
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
        </div>
      </div>

      <div className="bg-white p-4 rounded">
        <Tabs defaultValue={activeTab} onValueChange={handleTabChange}>
          <TabsList className="w-full justify-start">
            <TabsTrigger value="queues">Queues</TabsTrigger>
            <TabsTrigger value="pubsub">PubSub</TabsTrigger>
            <TabsTrigger value="users&access">Users & Access</TabsTrigger>
            <TabsTrigger value="settings">Settings</TabsTrigger>
          </TabsList>

          {breadcrumb && (
            <div className="py-4">
              <div className="flex items-center text-sm text-gray-600">
                <a href="/" className="hover:text-gray-900">Queues</a>
                <ChevronRight className="h-4 w-4 mx-2" />
                <span className="text-gray-900">{breadcrumb}</span>
              </div>
            </div>
          )}

          {!hideContent && (
            <>
              {activeTab === "queues" && <Queues />}
              <TabsContent value="pubsub">Make changes to your account here.</TabsContent>
              <TabsContent value="users&access">Change your password here.</TabsContent>
              <TabsContent value="settings">Change your password here.</TabsContent>
            </>
          )}
        </Tabs>
      </div>
    </div>
  )
}