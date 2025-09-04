import { useEffect } from "react";
import { Tabs, TabsContent } from "@/components/ui/tabs";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { Alert, AlertDescription } from "@/components/ui/alert";

export default function QueueDetails({ queueDetails, error }) {
  const { toast } = useToast();

  useEffect(() => {
    if (error && error.type === "API_UNREACHABLE") {
      toast({
        title: "Connection Error",
        description: error.message || "Failed to reach the backend. Please check your connection and try again.",
        variant: "destructive",
      });
    }
  }, [error, toast]);

  if (error) {
    if (error.type === "API_UNREACHABLE") {
      return (
        <TabsContent value="overview" className="space-y-4">
          <Alert variant="destructive">
            <AlertDescription>
              {error.message || "Could not connect to the backend API. Please ensure the API is running and accessible, then refresh the page."}
            </AlertDescription>
          </Alert>
        </TabsContent>
      );
    }
    return (
      <TabsContent value="overview" className="space-y-4">
        <Alert variant="destructive">
            <AlertDescription>
                Error: {typeof error === 'string' ? error : error.message || "An unexpected error occurred."}
            </AlertDescription>
        </Alert>
      </TabsContent>
    );
  }

  if (!queueDetails) {
    return (
      <div className="p-4">
        <div>Loading...</div>
      </div>
    );
  }

  return (
    <TabsContent value="overview" className="space-y-4">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Queue Name</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{queueDetails.queueName}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Queue ID</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{queueDetails.queueId}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Created At</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {new Date(queueDetails.createdAt).toLocaleString()}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Max Receive Attempts</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{queueDetails.maxReceiveAttempts}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Retention Period</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{queueDetails.retentionPeriodSeconds}s</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Visibility Timeout</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{queueDetails.visibilityTimeoutSeconds}s</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Eviction Policy</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {
                {
                  undefined: "Unspecified",
                  EVICTION_POLICY_UNSPECIFIED: "Unspecified",
                  EVICTION_POLICY_DROP: "Drop Message",
                  EVICTION_POLICY_DEAD_LETTER: "Dead Letter Queue",
                  EVICTION_POLICY_REORDER: "Reorder Message",
                }[queueDetails.evictionPolicy]
              }
            </div>
          </CardContent>
        </Card>
      </div>
    </TabsContent>
  );
}
