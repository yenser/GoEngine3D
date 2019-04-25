package main

import "github.com/go-gl/gl/v4.1-core/gl"

func buildVAO() {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
}

// struct Object {
// 	GLuint buffers[2]; // first provides position, second provides normal
// 	size_t offsets[2]; // beginning of vertex data in buffer
// 	size_t strides[2]; // stride of entire element in buffer

// 	unsigned int triangles;
// };

// void draw(unsigned int count, const struct Object *objects) {
// 	glBindVertexArray(vao);

// 	for(int i = 0; i < count; ++i) {
// 			// consecutive vertex buffer binding points are used here
// 			glBindVertexBuffers(0, 2, objects[i].buffers, objects[i].offsets, objects[i].strides);
// 			glDrawArrays(GL_TRIANGLES, 0, objects[i].triangles*3);
// 	}

// 	glBindVertexArray(0);
// }
// 			// consecutive vertex buffer binding points are used here
// 			glBindVertexBuffers(0, 2, objects[i].buffers, objects[i].offsets, objects[i].strides);
// 			glDrawArrays(GL_TRIANGLES, 0, objects[i].triangles*3);
// 	}

// 	glBindVertexArray(0);
// }
// 			// consecutive vertex buffer binding points are used here
// 			glBindVertexBuffers(0, 2, objects[i].buffers, objects[i].offsets, objects[i].strides);
// 			glDrawArrays(GL_TRIANGLES, 0, objects[i].triangles*3);
// 	}

// 	glBindVertexArray(0);
// }
