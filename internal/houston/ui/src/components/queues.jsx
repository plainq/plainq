import { useEffect, useState } from "react";
import { TabsContent } from "@/components/ui/tabs";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { QueueCreateDialog } from "./queueCreateDialog";
import QueuePagination from "./queuePagination";
import { useToast } from "@/hooks/use-toast";
import { Alert, AlertDescription } from "@/components/ui/alert";

export default function Queues() {  
  const { toast } = useToast();
  const [queues, setQueues] = useState([]);
  const [loading, setLoading] = useState(true);
  const [cursor, setNextCursor] = useState("");
  const [hasMore, setHasMore] = useState(false);
  const [limit, setLimit] = useState(10);
  const [cursors, setCursors] = useState(['']);
  const [currentPage, setCurrentPage] = useState(0);
  const [isApiUnreachable, setIsApiUnreachable] = useState(false);
  
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [queueName, setQueueName] = useState("");
  const [retentionPeriod, setRetentionPeriod] = useState(86400);
  const [visibilityTimeout, setVisibilityTimeout] = useState(30);
  const [receiveAttempts, setReceiveAttempts] = useState(5);
  const [dropPolicy, setDropPolicy] = useState("EVICTION_POLICY_DROP");
  const [deadLetterQueueId, setDeadLetterQueueId] = useState("");

  useEffect(() => {
    // Load stored values when component mounts
    const storedCursors = localStorage.getItem('queue_cursors');
    const storedPage = localStorage.getItem('queue_current_page');
    
    if (storedCursors) {
      setCursors(JSON.parse(storedCursors));
    }
    if (storedPage) {
      setCurrentPage(parseInt(storedPage));
    }
  }, []);

  useEffect(() => {
    fetchQueues(cursors[currentPage]);
  }, [limit, currentPage]);

  useEffect(() => {
    // Only save to localStorage on the client side
    if (typeof window !== 'undefined') {
      localStorage.setItem('queue_cursors', JSON.stringify(cursors));
      localStorage.setItem('queue_current_page', currentPage.toString());
    }
  }, [cursors, currentPage]);

  const handlePageChange = (direction) => {
    if (direction === 'next' && hasMore) {
      // Store the new cursor and move to next page
      setCursors(prev => {
        const newCursors = prev.slice(0, currentPage + 1);
        newCursors.push(cursor);
        return newCursors;
      });
      setCurrentPage(prev => prev + 1);
    } else if (direction === 'prev' && currentPage > 0) {
      setCurrentPage(prev => prev - 1);
    } else if (typeof direction === 'number' && direction >= 0 && direction <= cursors.length - 1) {
      setCurrentPage(direction);
    }
  };

  const handleLimitChange = (newLimit) => {
    setLimit(newLimit);
    setNextCursor("");
    setHasMore(false);
    setCursors(['']);
    setCurrentPage(0);
    fetchQueues("");
  };

  const fetchQueues = async (currentCursor = "") => {
    try {
      setLoading(true);
      setIsApiUnreachable(false);
      const response = await fetch(
        `http://localhost:8081/api/v1/queue?limit=${limit}&cursor=${currentCursor}`
      );
      if (!response.ok) {
        let errorMsg = "Failed to fetch queues";
        try {
            const errorData = await response.json();
            if (errorData && errorData.message) {
                errorMsg = errorData.message;
            }
        } catch (e) {
            // Ignore if response is not JSON or already handled
        }
        throw new Error(errorMsg);
      }

      const data = await response.json();
      setQueues(data.queues || []);
      setNextCursor(data.nextCursor || "");
      setHasMore(data.hasMore || false);
    } catch (err) {
      console.error("Error fetching queues:", err);
      if (err instanceof TypeError) {
        toast({
          title: "Connection Error",
          description: "Failed to reach the backend. Please check your connection and try again.",
          variant: "destructive",
        });
        setIsApiUnreachable(true);
      } else {
        toast({
          title: "Error Fetching Queues",
          description: err.message || "An unexpected error occurred while fetching queues.",
          variant: "destructive",
        });
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCreateQueue = async () => {
    // Validate required fields
    if (!queueName?.trim()) {
      toast({
        title: "Validation Error",
        description: "Queue name is required.",
        variant: "destructive",
      });
      return;
    }

    try {
      setIsApiUnreachable(false);
      const response = await fetch("http://localhost:8081/api/v1/queue", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          queueName,
          retentionPeriodSeconds: parseInt(retentionPeriod),
          visibilityTimeoutSeconds: parseInt(visibilityTimeout),
          maxReceiveAttempts: parseInt(receiveAttempts),
          evictionPolicy: dropPolicy,
          deadLetterQueueId: deadLetterQueueId || undefined,
        }),
      });

      let responseData;
      
      try {
        responseData = await response.json();
      } catch (parseError) {
        throw new Error(response.statusText || "Server error occurred");
      }

      if (!response.ok) {
        throw new Error(responseData.message || "Failed to create queue");
      }

      // Reset form
      setQueueName("");
      setRetentionPeriod("86400");
      setVisibilityTimeout("30");
      setReceiveAttempts("5");
      setDropPolicy("EVICTION_POLICY_DROP");
      setDeadLetterQueueId("");
      
      // Close dialog and show success message
      setIsDialogOpen(false);
      toast({
        title: "Success",
        description: "Queue created successfully!",
      });
      window.location.href = `/queue/${responseData.queueId}`;
    } catch (err) {
      console.error("Error creating queue:", err);
      if (err instanceof TypeError) {
        toast({
          title: "Connection Error",
          description: "Failed to reach the backend. Please check your connection and try again.",
          variant: "destructive",
        });
        setIsApiUnreachable(true);
      } else {
        toast({
          title: "Error Creating Queue",
          description: err.message || "Failed to create queue.",
          variant: "destructive",
        });
      }
    }
  };

  if (loading) return <div>Loading...</div>;

  if (isApiUnreachable) {
    return (
      <div className="container mx-auto pt-3 px-2">
        <Alert variant="destructive">
          <AlertDescription>
            Could not connect to the backend API. Please ensure the API is running and accessible, then refresh the page.
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className="bg-white px-2">
      <TabsContent value="queues">
        <div>
          <div className="flex flex-row justify-between pt-4 pb-4">
            <div>
              <p className="text-2xl font-bold">Queues</p>
            </div>

            <div className="flex flex-row gap-2 justify-end">
              <Button onClick={() => setIsDialogOpen(true)}>Create Queue</Button>
            </div>
          </div>

          <QueueCreateDialog
            isOpen={isDialogOpen}
            onOpenChange={setIsDialogOpen}
            onCreateQueue={handleCreateQueue}
            queueName={queueName}
            setQueueName={setQueueName}
            retentionPeriod={retentionPeriod}
            setRetentionPeriod={setRetentionPeriod}
            visibilityTimeout={visibilityTimeout}
            setVisibilityTimeout={setVisibilityTimeout}
            receiveAttempts={receiveAttempts}
            setReceiveAttempts={setReceiveAttempts}
            dropPolicy={dropPolicy}
            setDropPolicy={setDropPolicy}
            deadLetterQueueId={deadLetterQueueId}
            setDeadLetterQueueId={setDeadLetterQueueId}
          />

          <Table className="w-full">
            {/* <TableCaption>List of queues</TableCaption> */}
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>ID</TableHead>
                <TableHead>Created At</TableHead>
                <TableHead>Attempts</TableHead>
                <TableHead>Retention Period</TableHead>
                <TableHead>Visibility Timeout</TableHead>
                <TableHead>Eviction Policy</TableHead>
                <TableHead>Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {queues.map((queue) => {
                console.log("Rendering queue:", queue);
                return (
                  <TableRow key={queue.queueId}>
                    <TableCell>
                      <a href={`/queue/${queue.queueId}`}>{queue.queueName}</a>
                    </TableCell>
                    <TableCell>{queue.queueId}</TableCell>
                    <TableCell>
                      {new Date(queue.createdAt).toLocaleString("en-US", {
                        year: "numeric",
                        month: "2-digit",
                        day: "2-digit",
                        hour: "2-digit",
                        minute: "2-digit",
                        second: "2-digit",
                        hour12: false,
                      })}
                    </TableCell>
                    <TableCell>{queue.maxReceiveAttempts}</TableCell>
                    <TableCell>{queue.retentionPeriodSeconds}</TableCell>
                    <TableCell>{queue.visibilityTimeoutSeconds}</TableCell>
                    <TableCell>
                      {
                        {
                          undefined: "Unspecified",
                          EVICTION_POLICY_UNSPECIFIED: "Unspecified",
                          EVICTION_POLICY_DROP: "Drop Message",
                          EVICTION_POLICY_DEAD_LETTER: "Dead Letter Queue",
                          EVICTION_POLICY_REORDER: "Reorder Message",
                        }[queue.evictionPolicy]
                      }
                    </TableCell>
                    <TableCell>
                      <Button variant="outline">...</Button>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
          
          <QueuePagination
            hasMore={hasMore}
            onPageChange={handlePageChange}
            limit={limit}
            onLimitChange={handleLimitChange}
            currentPage={currentPage}
            totalPages={cursors.length}
          />
        </div>
      </TabsContent>
    </div>
  );
}
