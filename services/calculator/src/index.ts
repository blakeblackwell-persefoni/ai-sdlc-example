import { Hono } from "hono";
import { calculator } from "./routes/calculator";

const app = new Hono();

app.route("/", calculator);

app.notFound((c) => {
  return c.json({ error: "Not found" }, 404);
});

app.onError((err, c) => {
  console.error("Unhandled error:", err);
  return c.json({ error: "Internal server error" }, 500);
});

export default app;
