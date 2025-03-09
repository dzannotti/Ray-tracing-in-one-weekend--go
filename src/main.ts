const canvas = document.querySelector("canvas");
if (!canvas) throw new Error("No canvas");

const width = 512;
const aspectRatio = 1;
const height = width / aspectRatio;
canvas.width = width;
canvas.height = height;

const ctx = canvas?.getContext("2d");
if (!ctx) throw new Error("No ctx");

function fillRGB(r: number, g: number, b: number) {
  ctx!.fillStyle = `rgb(${r},${g},${b})`;
}

for (let j = 0; j <= height - 1; j++) {
  console.log(`Rendering scanline ${j}`);
  for (let i = 0; i <= width - 1; i++) {
    const r = Math.floor((i / width) * 255);
    const g = Math.floor((j / height) * 255);

    fillRGB(r, g, 0);
    ctx.fillRect(i, j, 1, 1);
  }
}
