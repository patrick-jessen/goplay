#version 330 core
layout (location = 0) out vec3 fragCol;

in vec3 fragPos;
in vec3 fragNorm;
in vec2 fragUV;
in mat3 TBN;

uniform sampler2D tex0; // Diffuse
uniform sampler2D tex1; // Normal

layout (std140) uniform shader_data {
  mat4 viewProjMat;
  mat4 modelMat;
  vec4 viewPos;
};


const vec4 ambient = vec4(1, 1, 1, 1.0) * 0.1;
const float specularStrength = 1.8;
const int specularPower = 32;

struct DirLight {
  vec3 direction;
  vec3 color;
};
DirLight dirLight = DirLight(normalize(vec3(0, 0, -5)), vec3(1,1,1));


////////////////////////////////////////////////////////////////////////////////
vec3 calcNormal() {
  vec3 n = texture(tex1, fragUV).rgb;
  n = normalize(n * 2 - 1);
  return normalize(TBN * n);
}

////////////////////////////////////////////////////////////////////////////////
vec4 calcDirectionalLight(DirLight light, vec3 normal) {
  // Diffuse
  vec3 lightVec = normalize(-light.direction);
  float difStrength = max(dot(normal, lightVec), 0.0);
  vec3 diffuse = difStrength * light.color;

  // Specular
  vec3 viewDir = normalize(vec3(viewPos) - fragPos);
  vec3 reflectDir = reflect(-lightVec, normal);  
  float specStrength = pow(max(dot(viewDir, reflectDir), 0.0), specularPower);
  vec3 specular = specularStrength * specStrength * light.color;
  
  return vec4(diffuse + specular, 1.0);
}

////////////////////////////////////////////////////////////////////////////////
void main() {
  vec3 normal = calcNormal();
  vec4 lights = calcDirectionalLight(dirLight, normal);

  vec4 texCol = texture(tex0, fragUV);
  fragCol = vec3(texCol * (ambient + lights));
}