-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    parent_id UUID REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create product_categories table
CREATE TABLE IF NOT EXISTS product_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert some sample categories
INSERT INTO categories (name, description) VALUES 
('Science', 'Science-related topics including biology, chemistry, physics'),
('Mathematics', 'Mathematics topics including algebra, calculus, statistics'),
('History', 'Historical topics and events'),
('Language', 'Language learning and literature'),
('Technology', 'Computer science and technology topics'),
('Education', 'General educational topics'),
('Biology', 'Study of living organisms'),
('Chemistry', 'Study of matter and chemical reactions'),
('Physics', 'Study of matter, energy, and their interactions'),
('Geography', 'Study of Earth''s physical features');

-- Set some parent-child relationships
UPDATE categories SET parent_id = (SELECT id FROM categories WHERE name = 'Science') WHERE name IN ('Biology', 'Chemistry', 'Physics');
UPDATE categories SET parent_id = (SELECT id FROM categories WHERE name = 'Education') WHERE name IN ('Science', 'Mathematics', 'History', 'Language', 'Technology');
