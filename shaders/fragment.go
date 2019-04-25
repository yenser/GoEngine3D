package shaders

var FragmentShader = `
#version 410 core

out vec4 outColor;

uniform sampler2D tex;

void main()
{
    outColor = vec4(1.0, 1.0, 1.0, 1.0);
}
` + "\x00"

// var FragmentShader = `
// #version 410 core
// in vec3 Color;
// in vec2 Texcoord;

// out vec4 outColor;

// uniform sampler2D tex;

// void main()
// {
//     outColor = texture(tex, Texcoord) * vec4(Color, 1.0);
// }
// ` + "\x00"
