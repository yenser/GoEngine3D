package shaders

var VertexShader = `
#version 410 core

in vec3 position;

out vec3 Color;

uniform mat4 model;
uniform mat4 camera;
uniform mat4 projection;
uniform vec3 overrideColor;

void main()
{
    gl_Position = projection * camera * model * vec4(position, 1.0);
}
` + "\x00"
