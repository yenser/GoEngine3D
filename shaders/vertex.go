package shaders

var VertexShader = `
#version 410 core

in vec3 position;

out vec3 Color;
out vec3 FragPos;
out vec3 Normal;

uniform mat4 model;
uniform mat4 camera;
uniform mat4 projection;

void main()
{
    gl_Position = projection * camera * model * vec4(position, 1.0);
    FragPos = vec3(model * vec4(1.0));
    Normal = normalize(model);
}
` + "\x00"
