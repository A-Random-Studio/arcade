<template>
    <div>
        <canvas ref="canvas" :width="width" :height="height"></canvas>
    </div>
</template>

<script>
import { EventBus } from "@/eventBus.js";
import { DrawEvent } from "@/modules/scribble/utility/drawEvents";

export default {
    name: "Canvas",

    props: {
        width: Number,
        height: Number,
        defaultBrushStyle: Object,
        drawingLocked: Boolean,
    },

    data: function () {
        return {
            mouseDown: false,
            previousPosition: { x: 0, y: 0 },
            brushStyle: this.defaultBrushStyle,
            canvas: null,
            context: null,
            offsetLeft: 0,
            offsetTop: 0,
        };
    },

    mounted: function () {
        this.canvas = this.$refs["canvas"];
        this.context = this.canvas.getContext("2d");
        this.canvas.addEventListener("mousemove", this.onMouseMove, false);
        this.canvas.addEventListener("mousedown", this.onMouseDown, false);
        this.canvas.addEventListener("mouseup", this.onMouseUp, false);
        this.canvas.addEventListener("mouseover", this.onMouseOver, false);
        EventBus.$on(DrawEvent.UPDATE_BRUSH, this.setBrushStyle);
    },

    methods: {
        onMouseDown: function (event) {
            const rect = this.canvas.getBoundingClientRect();

            this.previousPosition = {
                x:
                    (event.clientX - rect.left) *
                    (this.canvas.width / rect.width),
                y:
                    (event.clientY - rect.top) *
                    (this.canvas.height / rect.height),
            };

            this.handleDrawInput(
                this.previousPosition,
                this.previousPosition,
                this.brushStyle
            );
            this.mouseDown = true;
        },

        onMouseMove: function (event) {
            if (this.mouseDown) {
                const rect = this.canvas.getBoundingClientRect();

                const currentPosition = {
                    x:
                        (event.clientX - rect.left) *
                        (this.canvas.width / rect.width),
                    y:
                        (event.clientY - rect.top) *
                        (this.canvas.height / rect.height),
                };
                this.handleDrawInput(
                    this.previousPosition,
                    currentPosition,
                    this.brushStyle
                );
                this.previousPosition = currentPosition;
            }
        },

        onMouseUp: function () {
            this.mouseDown = false;
        },

        onMouseOver: function (event) {
            if (
                !(event.buttons === undefined
                    ? (event.which & 1) === 1
                    : (event.buttons & 1) === 1)
            ) {
                this.mouseDown = false;
            } else {
                this.mouseDown = true;
                const rect = this.canvas.getBoundingClientRect();

                this.previousPosition = {
                    x:
                        (event.clientX - rect.left) *
                        (this.canvas.width / rect.width),
                    y:
                        (event.clientY - rect.top) *
                        (this.canvas.height / rect.height),
                };
            }
        },

        setBrushStyle: function (brushStyle) {
            this.brushStyle = brushStyle;
        },

        handleDrawInput: function (from, to, brushStyle) {
            if (!this.drawingLocked) {
                const drawAction = {
                    from: from,
                    to: to,
                    brushStyle: brushStyle,
                    lineCap: this.context.lineCap,
                };
                this.draw(drawAction);
                this.$emit("drawAction", drawAction);
            }
        },

        draw: function (drawAction) {
            this.context.beginPath();
            this.context.moveTo(drawAction.from.x, drawAction.from.y);
            this.context.lineTo(drawAction.to.x, drawAction.to.y);
            this.context.strokeStyle = drawAction.brushStyle.brushColor;
            this.context.lineWidth = drawAction.brushStyle.brushSize;
            this.context.lineCap = "round";
            this.context.stroke();
            this.context.closePath();
        },

        resetCanvas: function() {
            this.context.clearRect(0, 0, this.width, this.height);
            // TODO: Add a way to reset canvas on the api
        }
    },
};
</script>

<style scoped>
canvas {
    top: 10%;
    left: 10%;
    border: 2px solid;
    width: 100%;
}
</style>
