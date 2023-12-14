# Presio Backgrounds

A small image generator written in GOLANG.

![test](./generator/out.png)

# Quickstart local development

- Rename `env.localhost` file to be picked up by docker compose. It has defaults values for local development

```bash
mv ./.env.localhost ./.env
```

- Start all containers via `docker compose`

```bash
docker compose up -d
```

- Generate an image using a post request with one of the `_exemples` body

```bash
curl -X POST 'http://localhost:4444/image'
   -H 'Content-Type: application/json'
   -d '${the_json_body}'
```

- Profit from your generated (base 64 if desired) image !

# API Reference

```typescript
interface Pattern {
  pattern_type: "PLAIN_COLOR" | "LINEAR_GRADIENT";
  colors: string[]; // Array of hex colors
  angle: number; // Int 0-360, used for gradients
}

interface Blob {
  shape_type: "CIRCLE" | "BLOB" | "POLYGON";

  /* Array of coordinates within a -100;100 space
    When `shape_type` is `CIRCLE`, only one point is needed, representing the radius of the circle
  */
  coordinates: number[][];
  /* Coordinates, between 0-1. Relative to the size of the whole drawing area */
  center: number[];
  /* Pattern to draw the shape with */
  pattern: Pattern;
  /* Rotation in degrees, from 0 to 360 */
  rotation: number;
  /* float64. Stroke width, if any. By default, fill the shape instead  */
  stroke_width: number;
  /* Scaling: 1 = 100% */
  scale: number;
  /* Elevation, unitless but still a positive integer */
  elevation: number;
}

interface SteppedDrawingOptions {
  /* How many steps will be executed. Integer */
  steps: number;
  /* The difference in scale between steps. Float */
  scale_step: number;
  /* The difference in rotation between steps. Integer */
  rotation_step: number;
  /* The difference in elevation between steps. Positive integer */
  elevation_step: number;
  /* The difference in coordinates between steps
   * Note that the coordinates interpretation will differ according to Blob settings
   * Array of 2 floats (0-1)
   */
  translate_step: number[];
  /* Start and end patterns, lerping from one to another */
  start_pattern: Pattern;
  end_pattern: Pattern;
}

interface DrawingInstruction {
  /* The blob to draw */
  blob: Blob;
  /* If desired, draw the shape repeatedly */
  options: SteppedDrawingOptions;
}

interface PostCreateImage {
  /* Positive integers */
  width: number;
  height: number;

  /* Background of the image */
  background: Pattern;
  /* The list of shapes to draw */
  blobs: DrawingInstruction[];
}
```

# TODOs

- Avoid init one connection per API call...
- Optimize use of drawing layers to save GC time
