package shaders

var FragmentShader = `
#version 410 core

in vec3 Color;
in vec3 FragPos;
in vec3 Normal;
out vec4 outColor;

uniform sampler2D tex;
uniform vec3 objectColor;
uniform vec3 lightColor;

void main()
{

    vec3 lightPos = vec3(0, 1, 0);
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);

    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;

    float ambientStrength = 0.2;
    vec3 ambient = ambientStrength * lightColor;

    vec3 result = (ambient + diffuse) * objectColor;

    outColor = vec4(result, 1.0);
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
