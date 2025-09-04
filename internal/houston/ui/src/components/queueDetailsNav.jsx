import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs"
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar"
import QueueDetails from "@/components/queueDetails.jsx";
import { ChevronRight } from "lucide-react";

export default function QueueDetailsNav({queueDetails, error, breadcrumbText}) {
  return (
    <div className="bg-white p-4 rounded">
      {breadcrumbText && (
        <div className="pb-4 mb-4 border-b border-gray-200">
          <div className="flex items-center text-sm text-gray-600">
            <a href="/" className="hover:text-gray-900">Queues</a>
            <ChevronRight className="h-4 w-4 mx-2 shrink-0" />
            <span className="font-medium text-gray-900 truncate" title={breadcrumbText}>{breadcrumbText}</span>
          </div>
        </div>
      )}
      <Tabs defaultValue="overview">
        <TabsList className="w-full justify-start">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="messages">Messages</TabsTrigger>
          <TabsTrigger value="metrics">Metrics</TabsTrigger>
          <TabsTrigger value="settings">Settings</TabsTrigger>
        </TabsList>
        <QueueDetails queueDetails={queueDetails} error={error}/>
        <TabsContent value="messages" className="px-1">Queue messages will be shown here.</TabsContent>
        <TabsContent value="metrics" className="px-1">Queue metrics will be shown here.</TabsContent>
        <TabsContent value="settings" className="px-1">Queue settings will be shown here.</TabsContent>
      </Tabs>
    </div>
  )
}